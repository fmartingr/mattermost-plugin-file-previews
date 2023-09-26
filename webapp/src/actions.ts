import {GlobalState} from 'mattermost-redux/types/store';
import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {PluginId} from './plugin_id';

export const getSiteURL = (state: GlobalState): string => {
    const config = getConfig(state);

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substring(0, basePath.length - 1);
        }
    }

    return basePath;
};

export const getPluginServerRoute = (state: GlobalState): string => {
    const siteURL = getSiteURL(state);
    return `${siteURL}/plugins/${PluginId}`;
};
