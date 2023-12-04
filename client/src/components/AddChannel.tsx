import { useState } from 'react'
import { useForm } from '@mantine/form';
import { Button, Modal, Group, TextInput, Textarea, Select } from '@mantine/core'
import { ENDPOINT, YTChannel } from "../App"
import { KeyedMutator } from 'swr';
import SearchBar from './SearchBar';

function AddChannel({
    setSearchedChannels,
}: {
    setSearchedChannels: React.Dispatch<React.SetStateAction<YTChannel[]>>;
}) {
    const [query, setQuery] = useState('')
    const [platform, setPlatform] = useState('');

    const handleSearch = (searchText: string) => {
        setQuery(searchText);
    };

    async function handleSearchEnter(title: string) {
        const data = await fetch(`${ENDPOINT}/stats/${platform}/${title}`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json"
            },
        }).then((r) => r.json());

        setSearchedChannels(data)
    }

    return (
        <>
            <SearchBar
                onSearch={handleSearch}
                className="search-bar"
                onKeyDown={handleSearch}
                handleSearchEnter={handleSearchEnter}
            ></SearchBar>
            <Select
                placeholder="Streaming platform"
                data={['YouTube', 'Twitch', 'Kick']}
                value={platform}
                onChange={(value) => {
                    if (value !== null) {
                        setPlatform(value);
                    }
                }}
                clearable
                style={{ width: '300px' }}
            ></Select>
            <Button type="submit" onClick={() => handleSearchEnter(query)}>Add channel</Button>
        </>
    )
}

export default AddChannel;