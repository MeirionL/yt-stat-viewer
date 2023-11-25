import SearchBar from "../client/src/components/SearchBar";
import { useState, useEffect } from "react";
import channelsData from "../client/src/data/channels";
import Dropdown from "../client/src/components/Dropdown";
import './App.css';




function App() {
    const [query, setQuery] = useState<string>('')
    const [selectedPlatform, setSelectedPlatform] = useState<string>('');
    const channels = channelsData;

    const handleSearch = (searchText: string) => {
        setQuery(searchText);
    };

    const filteredChannels = channels.filter((channel) => {
        const queryLower = query.toLowerCase();
        const channelNameLower = channel.name.toLowerCase();

        return (selectedPlatform === '' || channel.platform.toLowerCase() === selectedPlatform.toLowerCase()) &&
            channelNameLower.startsWith(queryLower);
    });

    const handlePlatformChange = (newPlatform: string) => {
        setSelectedPlatform(newPlatform);
    }

    return (
        <div>
            <SearchBar onSearch={handleSearch} className="search-bar" />
            <br></br>
            <br></br>
            <Dropdown onClick={handlePlatformChange} selectedPlatform={selectedPlatform}></Dropdown>
            <ul>
                {filteredChannels.map((channel, index) => (
                    <li key={index}>
                        {channel.name}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default App;