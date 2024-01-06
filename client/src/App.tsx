import { Box, List, ThemeIcon, MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import { CheckCircleFillIcon } from "@primer/octicons-react";
import { useState } from 'react';
import AddChannel from './components/AddChannel';
import GrantPermissions from "./components/GrantPermissions";

export interface YTChannel {
    title: string;
    subscribers: number;
    videos: number;
    views: number;
    last_stream_time: string;
    is_live: string;
    stream_title: string;
}

export const ENDPOINT = 'http://localhost:8080';

function App() {

    const [displayedChannel, setDisplayedChannel] = useState<YTChannel | null>(null);
    const [searchedChannels, setSearchedChannels] = useState<YTChannel[]>([]);

    function showChannelStats() {
        if (displayedChannel) {
            let streamTitleLabel;
            if (displayedChannel.is_live === "No") {
                streamTitleLabel = 'Previous stream title';
            } else {
                streamTitleLabel = 'Current stream title';
            }

            return (
                <div>
                    <h2>{displayedChannel.title}</h2>
                    <p>Subscribers: {displayedChannel.subscribers}</p>
                    <p>Videos: {displayedChannel.videos}</p>
                    <p>Views: {displayedChannel.views}</p>
                    {displayedChannel.last_stream_time !== "" && (
                        <>
                            <p>Start time of last stream: {displayedChannel.last_stream_time}</p>
                            <p>Is currently live: {displayedChannel.is_live}</p>
                            <p>{streamTitleLabel}: {displayedChannel.stream_title}</p>
                        </>
                    )}
                </div>
            );
        }
        return (
            <div>
                <h2>No channel selected</h2>
            </div>
        );
    }

    return (
        <MantineProvider>
            <Box m="xl" p="xl">
                <List spacing="xl" size="xl" mb={12} center>
                    {searchedChannels?.map((channel: YTChannel) => {
                        return (
                            <List.Item
                                maw="xl"
                                onClick={() => setDisplayedChannel(channel)}
                                key={`channels_list__${channel.title}`}
                                icon={
                                    <ThemeIcon color="red" size={24} radius="xl">
                                        <CheckCircleFillIcon size={20} />
                                    </ThemeIcon>
                                }
                            >
                                {channel.title}
                            </List.Item>
                        );
                    })}
                </List>
                <AddChannel setSearchedChannels={setSearchedChannels} />
                {showChannelStats()}
                <GrantPermissions />
            </Box>
        </MantineProvider>
    );
}

export default App;
