import { useState } from 'react';
import { Button } from '@mantine/core';
import { ENDPOINT, YTChannel } from '../App';
import SearchBar from './SearchBar';

interface AddChannelProps {
    setSearchedChannels: React.Dispatch<React.SetStateAction<YTChannel[]>>;
}

function AddChannel({ setSearchedChannels }: AddChannelProps) {
    const [query, setQuery] = useState('');

    const handleSearch = (searchText: string): void => {
        setQuery(searchText);
    };

    const handleSearchEnter = async (title: string): Promise<void> => {
        const data: YTChannel[] = await fetch(`${ENDPOINT}/stats/${title}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        }).then((r) => r.json());

        setSearchedChannels(data);
    };

    return (
        <>
            <SearchBar onSearch={handleSearch} handleSearchEnter={handleSearchEnter} />
            <Button onClick={() => handleSearchEnter(query)}>Add channel</Button>
        </>
    );
}

export default AddChannel;
