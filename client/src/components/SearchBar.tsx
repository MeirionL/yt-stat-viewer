import { ChangeEvent, KeyboardEvent } from 'react';

interface SearchBarProps {
    onSearch: (searchText: string) => void;
    handleSearchEnter: (title: string) => Promise<void>;
}

const SearchBar = ({ onSearch, handleSearchEnter }: SearchBarProps) => {
    const handleInputChange = (e: ChangeEvent<HTMLInputElement>): void => {
        onSearch(e.target.value);
    };

    const handleKeyPress = async (e: KeyboardEvent<HTMLInputElement>): Promise<void> => {
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
