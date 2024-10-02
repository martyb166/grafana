package resource

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
)

type IndexServer struct {
	ResourceServer
	s     *server
	index *Index
	ws    *indexWatchServer
}

func (is IndexServer) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	res := &SearchResponse{}
	return res, nil
}

func (is IndexServer) History(ctx context.Context, req *HistoryRequest) (*HistoryResponse, error) {
	return nil, nil
}

func (is IndexServer) Origin(ctx context.Context, req *OriginRequest) (*OriginResponse, error) {
	return nil, nil
}

// Load the index
func (is *IndexServer) Load(ctx context.Context) error {
	is.index = NewIndex(is.s, Opts{})
	return is.index.Init(ctx)
}

// Init sets the resource server on the index server
// so we can call the resource server from the index server
// TODO: a chicken and egg problem - index server needs the resource server but the resource server is created with the index server
func (is *IndexServer) Init(ctx context.Context, rs *server) error {
	is.s = rs
	is.ws = &indexWatchServer{
		is:      is,
		context: ctx,
	}
	// TODO: how to watch all resources?
	wr := &WatchRequest{
		Options: &ListOptions{
			Key: &ResourceKey{
				Group:    "playlist.grafana.app",
				Resource: "playlists",
			},
		},
	}
	// TODO: handle watch error
	var err error
	go func() {
		err = rs.Watch(wr, is.ws)
	}()
	return err
}

func NewResourceIndexServer() ResourceIndexServer {
	return &IndexServer{}
}

type ResourceIndexer interface {
	Index(ctx context.Context) (*Index, error)
}

type indexWatchServer struct {
	grpc.ServerStream
	context context.Context
	is      *IndexServer
}

func (f *indexWatchServer) Send(we *WatchEvent) error {
	if we.Type == WatchEvent_ADDED {
		return f.Add(we)
	}

	if we.Type == WatchEvent_DELETED {
		return f.Delete(we)
	}

	if we.Type == WatchEvent_MODIFIED {
		return f.Update(we)
	}

	return nil
}

func (f *indexWatchServer) RecvMsg(m interface{}) error {
	return nil
}

func (f *indexWatchServer) SendMsg(m interface{}) error {
	return errors.New("not implemented")
}

func (f *indexWatchServer) Context() context.Context {
	if f.context == nil {
		f.context = context.Background()
	}
	return f.context
}

func (f *indexWatchServer) Index() *Index {
	return f.is.index
}

func (f *indexWatchServer) Add(we *WatchEvent) error {
	data, err := getData(we.Resource)
	if err != nil {
		return err
	}
	err = f.Index().Index(f.context, data)
	if err != nil {
		return err
	}
	return nil
}

func (f *indexWatchServer) Delete(we *WatchEvent) error {
	// TODO: this seems flakey. Does a delete have a Resource or Previous?
	// both cases have happened ( maybe because Georges pr was reverted )
	rs := we.Resource
	if rs == nil {
		rs = we.Previous
	}
	if rs == nil {
		return errors.New("resource not found")
	}
	data, err := getData(rs)
	if err != nil {
		return err
	}
	err = f.Index().Delete(f.context, data.Uid, data.Key)
	if err != nil {
		return err
	}
	return nil
}

func (f *indexWatchServer) Update(we *WatchEvent) error {
	data, err := getData(we.Resource)
	if err != nil {
		return err
	}
	err = f.Index().Delete(f.context, data.Uid, data.Key)
	if err != nil {
		return err
	}
	err = f.Index().Index(f.context, data)
	if err != nil {
		return err
	}
	return nil
}

type Data struct {
	Key   *ResourceKey
	Value *ResourceWrapper
	Uid   string
}

func getGroup(r *Resource) string {
	v := strings.Split(r.ApiVersion, "/")
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

func getData(wr *WatchEvent_Resource) (*Data, error) {
	r, err := getResource(wr.Value)
	if err != nil {
		return nil, err
	}

	key := &ResourceKey{
		Group:     getGroup(r),
		Resource:  r.Kind,
		Namespace: r.Metadata.Namespace,
		Name:      r.Metadata.Name,
	}

	value := &ResourceWrapper{
		ResourceVersion: wr.Version,
		Value:           wr.Value,
	}
	return &Data{Key: key, Value: value, Uid: r.Metadata.Uid}, nil
}
