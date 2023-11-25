import { useState } from 'react'
import { useForm } from '@mantine/form';
import { Button, Modal, Group, TextInput, Textarea } from '@mantine/core'
import { ENDPOINT, YTChannel } from "../App"
import { KeyedMutator } from 'swr';
import SearchBar from './SearchBar';

function AddChannel({
    setSearchedChannels,
}: {
    setSearchedChannels: React.Dispatch<React.SetStateAction<YTChannel[]>>;
}) {
    const [query, setQuery] = useState<string>('')

    const handleSearch = (searchText: string) => {
        setQuery(searchText);
    };

    async function handleSearchEnter(title: string) {
        const data = await fetch(`${ENDPOINT}/youtube/stats/${title}`, {
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
            <Button type="submit" onClick={() => handleSearchEnter(query)}>Add channel</Button>
        </>
    )
}

export default AddChannel;