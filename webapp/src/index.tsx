import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/lib/store';

import {PluginId} from './plugin_id';

import {PluginRegistry} from '@/types/mattermost-webapp';

import CustomFilePreviewComponent from './components/preview';

import {FormattedMessage} from 'react-intl';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        registry.registerFilePreviewComponent(
            (fileInfo, post) => {
                if (fileInfo.extension === "docx") {
                    return true
                }
                return false
            },
            CustomFilePreviewComponent,
        );
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void
    }
}

window.registerPlugin(PluginId, new Plugin());
