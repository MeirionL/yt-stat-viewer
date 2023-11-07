import { ChangeEvent } from 'react';

interface SearchBarProps {
    onSearch: (searchText: string) => void;
    className?: string;
}

const SearchBar = ({ onSearch }: SearchBarProps) => {
    const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
        onSearch(e.target.value);
    };

    return (
        <>
            <input
                type="search"
                placeholder="Search..."
                onChange={handleInputChange}
            />
        </>
    );
};

export default SearchBar