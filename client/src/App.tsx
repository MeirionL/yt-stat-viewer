import useSWR from "swr"
import { Box, List, ThemeIcon, MantineProvider } from '@mantine/core'
import '@mantine/core/styles.css';
import { CheckCircleFillIcon } from "@primer/octicons-react";
import { useState } from 'react'
import AddChannel from './components/AddChannel';

export interface YTChannel {
    subscribers: number;
    title: string;
    views: number;
    platform: string;
}

export const ENDPOINT = 'http://localhost:8080'

const fetcher = (url: string) => fetch(`${ENDPOINT}/${url}`).then(r => r.json());

function getColorBasedOnTodoDone(platform: string | undefined): string {
    console.log(platform);
    if (platform === "YouTube") {
        return 'red';
    } else if (platform === 'twitch') {
        return 'purple';
    } else if (platform === 'kick') {
        return 'orange';
    }
    return 'grey';
}

function App() {

    const { data, mutate } = useSWR<YTChannel[]>(`youtube/stats`, fetcher);
    const [displayedChannel, setDisplayedChannel] = useState<YTChannel | null>(null);
    const [searchedChannels, setSearchedChannels] = useState<YTChannel[]>([]);

    function showChannelStats() {
        if (displayedChannel) {
            return (
                <div>
                    <h2>{displayedChannel.title}</h2>
                    <p>Subscribers: {displayedChannel.subscribers}</p>
                    <p>Views: {displayedChannel.views}</p>
                    <p>Platform: {displayedChannel.platform}</p>
                </div>
            );
        }
        return (
            <div>
                <h2>No channel selected</h2>
            </div>
        );
    }

    return <MantineProvider>{
        <Box m="xl" w="xl" p="xl">

            <List spacing="xl" size="xl" mb={12} center>
                {searchedChannels?.map((channel: YTChannel) => {
                    return (
                        <List.Item maw="xl"
                            onClick={() => setDisplayedChannel(channel)}
                            key={`channels_list__${channel.title}`}

                            icon={
                                <ThemeIcon color={getColorBasedOnTodoDone(channel.platform)} size={24} radius="xl">
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
        </Box>
    }</MantineProvider>;
}

export default App;