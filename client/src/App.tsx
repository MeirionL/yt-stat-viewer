import useSWR from "swr"
import { Box, List, ThemeIcon, MantineProvider } from '@mantine/core'
import '@mantine/core/styles.css';
import { CheckCircleFillIcon } from "@primer/octicons-react";
import { useState } from 'react'
import AddChannel from './components/AddChannel';
import GrantPermissions from "./components/GrantPermissions";

export interface YTChannel {
    subscribers: number;
    title: string;
    views: number;
    platform: string;
}

export const ENDPOINT = 'http://localhost:8080'

function App() {

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
        <Box m="xl" p="xl">

            <List spacing="xl" size="xl" mb={12} center>
                {searchedChannels?.map((channel: YTChannel) => {
                    return (
                        <List.Item maw="xl"
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
    }</MantineProvider>;
}

export default App;