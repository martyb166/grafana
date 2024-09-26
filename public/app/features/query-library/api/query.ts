import { BaseQueryFn } from '@reduxjs/toolkit/query/react';
import { lastValueFrom } from 'rxjs';

import { config } from '@grafana/runtime';
import { BackendSrvRequest, getBackendSrv, isFetchError } from '@grafana/runtime/src/services/backendSrv';

import { logQueryLibrary } from './logQueryLibrary';
import { DataQuerySpecResponse } from './types';

/**
 * @alpha
 */
export const API_VERSION = 'peakq.grafana.app/v0alpha1';

/**
 * @alpha
 */
export enum QueryTemplateKinds {
  QueryTemplate = 'QueryTemplate',
}

/**
 * Query Library is an experimental feature. API (including the URL path) will likely change.
 *
 * @alpha
 */
export const BASE_URL = `/apis/${API_VERSION}/namespaces/${config.namespace}/querytemplates/`;

// URL is optional for these requests
interface QueryLibraryBackendRequest extends Pick<BackendSrvRequest, 'data' | 'method'> {
  url?: string;
  headers?: { [key: string]: string };
}

/**
 * TODO: similar code is duplicated in many places. To be unified in #86960
 */
export const baseQuery: BaseQueryFn<QueryLibraryBackendRequest, DataQuerySpecResponse, Error> = async (
  requestOptions
) => {
  const start = performance.now();
  try {
    const responseObservable = getBackendSrv().fetch<DataQuerySpecResponse>({
      url: `${BASE_URL}${requestOptions.url ?? ''}`,
      showErrorAlert: true,
      method: requestOptions.method || 'GET',
      data: requestOptions.data,
      headers: { ...requestOptions.headers },
    });

    const response = await lastValueFrom(responseObservable);

    const recordCount = response.data.items.length || 0;

    const end = performance.now();
    const timeTaken = end - start;
    console.log({ type: 'backend fetch', recordCount, timeTaken: `${timeTaken} ms` });
    // log to a log file
    logQueryLibrary('backend fetch', recordCount, timeTaken);
    return response;
  } catch (error) {
    if (isFetchError(error)) {
      return { error: new Error(error.data.message) };
    } else if (error instanceof Error) {
      return { error };
    } else {
      return { error: new Error('Unknown error') };
    }
  }
};
