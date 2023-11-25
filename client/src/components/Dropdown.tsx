import { useState } from 'react';


interface DropdownProps {
    onClick: (platform: string) => void;
    selectedPlatform: string;
}

const Dropdown = ({ onClick, selectedPlatform }: DropdownProps) => {
    const [platform, setPlatform] = useState(selectedPlatform);

    const handleOptionClick = (option: string) => {
        setPlatform(option);
        onClick(option);
    };

    return (
        <>
            <div className="btn-group">
                <button type="button" className="btn btn-primary dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
                    {selectedPlatform === '' ? 'Platform' : selectedPlatform}
                </button>
                <ul className="dropdown-menu">
                    <li><a className="dropdown-item" href="#" onClick={() => handleOptionClick('Twitch')}>Twitch</a></li>
                    <li><a className="dropdown-item" href="#" onClick={() => handleOptionClick('YouTube')}>YouTube</a></li>
                    <li><a className="dropdown-item" href="#" onClick={() => handleOptionClick('Kick')}>Kick</a></li>
                </ul>
            </div>
        </>
    );
}

export default Dropdown;