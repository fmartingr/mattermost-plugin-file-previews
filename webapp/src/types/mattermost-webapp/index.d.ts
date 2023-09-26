import {FileInfo} from '@mattermost/types/lib/files';
import {Post} from '@mattermost/types/lib/posts';

export interface PluginRegistry {
    registerPostTypeComponent(typeName: string, component: React.ElementType)
    registerFilePreviewComponent(override: (fileInfo: FileInfo, post: Post) => bool, component: React.ElementType)
    registerPostDropdownMenuAction(text: string, action: func, filter: func)
    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}
