import { ChangeEvent, KeyboardEvent } from 'react';

interface SearchBarProps {
    onSearch: (searchText: string) => void;
    onKeyDown: (searchText: string) => void;
    handleSearchEnter: (title: string) => Promise<void>;
    className?: string;
}

const SearchBar = ({ onSearch, handleSearchEnter }: SearchBarProps) => {
    const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
        onSearch(e.target.value);
    };

    const handleKeyPress = async (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            await handleSearchEnter(e.currentTarget.value);
        }
    };

    return (
        <>
            <input
                type="search"
                placeholder="Search..."
                onChange={handleInputChange}
                onKeyDown={handleKeyPress}
            />
        </>
    );
};

export default SearchBar;