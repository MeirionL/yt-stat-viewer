import React from 'react';
import { YTChannel } from '../App';

interface ChannelInfoProps {
    channel: YTChannel;
}

const ChannelInfo: React.FC<ChannelInfoProps> = ({ channel }) => {
    const channelKeys = Object.keys(channel);

    return (
        <div>
            <h3>Attributes of {channel.title}:</h3>
            <ul>
                {channelKeys.map(key => (
                    <li key={key}>
                        <strong>{key}:</strong> {channel[key as keyof YTChannel]}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default ChannelInfo;
