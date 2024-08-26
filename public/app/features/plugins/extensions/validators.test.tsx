import { memo } from 'react';

import { PluginExtensionAddedLinkConfig, PluginExtensionLinkConfig, PluginExtensionTypes } from '@grafana/data';

import {
  assertConfigureIsValid,
  assertLinkPathIsValid,
  assertExtensionPointIdIsValid,
  assertStringProps,
  isReactComponent,
} from './validators';

describe('Plugin Extension Validators', () => {
  describe('assertLinkPathIsValid()', () => {
    it('should not throw an error if the link path is valid', () => {
      expect(() => {
        const pluginId = 'myorg-b-app';
        const extension = {
          path: `/a/${pluginId}/overview`,
          title: 'My Plugin',
          description: 'My Plugin Description',
          extensionPointId: '...',
        };

        assertLinkPathIsValid(pluginId, extension.path);
      }).not.toThrowError();
    });

    it('should throw an error if the link path is pointing to a different plugin', () => {
      expect(() => {
        const extension = {
          path: `/a/myorg-b-app/overview`,
          title: 'My Plugin',
          description: 'My Plugin Description',
          extensionPointId: '...',
        };

        assertLinkPathIsValid('another-plugin-app', extension.path);
      }).toThrowError();
    });

    it('should throw an error if the link path is not prefixed with "/a/<PLUGIN_ID>"', () => {
      expect(() => {
        const extension = {
          path: `/some-bad-path`,
          title: 'My Plugin',
          description: 'My Plugin Description',
          extensionPointId: '...',
        };

        assertLinkPathIsValid('myorg-b-app', extension.path);
      }).toThrowError();
    });
  });

  describe('assertExtensionPointIdIsValid()', () => {
    it('should throw an error if the extensionPointId does not have the right prefix', () => {
      expect(() => {
        assertExtensionPointIdIsValid('my-org-app', {
          type: PluginExtensionTypes.link,
          title: 'Title',
          description: 'Description',
          extensionPointId: 'wrong-extension-point-id',
        });
      }).toThrowError();
    });

    it('should NOT throw an error if the extensionPointId is correct', () => {
      expect(() => {
        assertExtensionPointIdIsValid('my-org-app', {
          type: PluginExtensionTypes.link,
          title: 'Title',
          description: 'Description',
          extensionPointId: 'grafana/some-page/extension-point-a',
        });

        assertExtensionPointIdIsValid('my-org-app', {
          type: PluginExtensionTypes.link,
          title: 'Title',
          description: 'Description',
          extensionPointId: 'plugins/my-super-plugin/some-page/extension-point-a',
        });
      }).not.toThrowError();
    });
  });

  describe('assertConfigureIsValid()', () => {
    it('should NOT throw an error if the configure() function is missing', () => {
      expect(() => {
        assertConfigureIsValid({
          title: 'Title',
          description: 'Description',
          targets: 'grafana/some-page/extension-point-a',
        } as PluginExtensionAddedLinkConfig);
      }).not.toThrowError();
    });

    it('should NOT throw an error if the configure() function is a valid function', () => {
      expect(() => {
        assertConfigureIsValid({
          title: 'Title',
          description: 'Description',
          targets: 'grafana/some-page/extension-point-a',
          configure: () => {},
        } as PluginExtensionAddedLinkConfig);
      }).not.toThrowError();
    });

    it('should throw an error if the configure() function is defined but is not a function', () => {
      expect(() => {
        assertConfigureIsValid(
          // @ts-ignore
          {
            title: 'Title',
            description: 'Description',
            extensionPointId: 'grafana/some-page/extension-point-a',
            handler: () => {},
            configure: '() => {}',
          } as PluginExtensionLinkConfig
        );
      }).toThrowError();
    });
  });

  describe('assertStringProps()', () => {
    it('should throw an error if any of the expected string properties is missing', () => {
      expect(() => {
        assertStringProps(
          {
            description: 'Description',
            extensionPointId: 'grafana/some-page/extension-point-a',
          },
          ['title', 'description', 'extensionPointId']
        );
      }).toThrowError();
    });

    it('should throw an error if any of the expected string properties is an empty string', () => {
      expect(() => {
        assertStringProps(
          {
            title: '',
            description: 'Description',
            extensionPointId: 'grafana/some-page/extension-point-a',
          },
          ['title', 'description', 'extensionPointId']
        );
      }).toThrowError();
    });

    it('should NOT throw an error if the expected string props are present and not empty', () => {
      expect(() => {
        assertStringProps(
          {
            title: 'Title',
            description: 'Description',
            extensionPointId: 'grafana/some-page/extension-point-a',
          },
          ['title', 'description', 'extensionPointId']
        );
      }).not.toThrowError();
    });

    it('should NOT throw an error if there are other existing and empty string properties, that we did not specify', () => {
      expect(() => {
        assertStringProps(
          {
            title: 'Title',
            description: 'Description',
            extensionPointId: 'grafana/some-page/extension-point-a',
            dontCare: '',
          },
          ['title', 'description', 'extensionPointId']
        );
      }).not.toThrowError();
    });
  });

  describe('isReactComponent()', () => {
    it('should return TRUE if we pass in a valid React component', () => {
      expect(isReactComponent(() => <div>Some text</div>)).toBe(true);
    });

    it('should return TRUE if we pass in a component wrapped with React.memo()', () => {
      const Component = () => <div>Some text</div>;
      const wrapped = memo(() => (
        <div>
          <Component />
        </div>
      ));
      wrapped.displayName = 'MyComponent';

      expect(isReactComponent(wrapped)).toBe(true);
    });

    it('should return FALSE if we pass in a valid React component', () => {
      expect(isReactComponent('Foo bar')).toBe(false);
      expect(isReactComponent(123)).toBe(false);
      expect(isReactComponent(false)).toBe(false);
      expect(isReactComponent(undefined)).toBe(false);
      expect(isReactComponent(null)).toBe(false);
    });
  });
});
