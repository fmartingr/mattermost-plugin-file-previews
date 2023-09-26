import React, { useMemo, useState } from 'react';

import Client4 from 'mattermost-redux/client/client4';

import { getPluginServerRoute } from '@/actions';
import { FileInfo } from '@mattermost/types/lib/files';
import { Post } from '@mattermost/types/lib/posts';
import { useDispatch } from 'react-redux';

const client = new Client4();

const getPluginFileURL = (fileId: string) => async (dispatch, getState): Promise<string> => {
    const state = getState();
    const baseURL = getPluginServerRoute(state);
    client.setUrl(baseURL);

    return Promise.resolve(`http://localhost:8065${baseURL}/preview/?fileID=${fileId}`)
};

type Props = {
    fileInfo: FileInfo
    post: Post
};

const CustomFilePreviewComponent = (props: Props) => {
    const dispatch = useDispatch();
    const [fileURL, setFileURL] = useState("");
    const className = "custom-file"

    useMemo(() => {
        // Dispatch the fetchData action creator and wait for the promise to resolve
        dispatch(getPluginFileURL(props.fileInfo.id))
            .then((url, fileInfo) => {
                setFileURL(url);
            })
            .catch((error) => {
                // Handle any errors here
                console.error('Error:', error);
            });
    }, [dispatch]);

    return (
        <iframe src={fileURL} width={"100%"} height={"100%"}></iframe>
    );
};

CustomFilePreviewComponent.propTypes = {
    fileInfo: PropTypes.object.isRequired,
    post: PropTypes.object.isRequired,
};

export default CustomFilePreviewComponent;
