import { useState } from 'react'
import { Button, Select } from '@mantine/core'

function GrantPermissions() {
    const [platform, setPlatform] = useState('');

    const handleGrantPermissions = () => {
        if (platform === 'YouTube') {
            window.location.href = "http://localhost:8080/auth/google";
        } else if (platform === 'Twitch') {
            window.location.href = "http://localhost:8080/auth/twitch";
        } else {
            console.log("Please select a platform")
        }
    };

    const handleRemovePermissions = () => {
        if (platform === 'YouTube') {
            window.location.href = "http://localhost:8080/logout/google";
        } else if (platform === 'Twitch') {
            window.location.href = "http://localhost:8080/logout/twitch";
        } else {
            console.log("Please select a platform")
        }
    };

    return (
        <>
            <Select
                placeholder="Permission streaming platform"
                data={['YouTube', 'Twitch']}
                value={platform}
                onChange={(value) => {
                    if (value !== null) {
                        setPlatform(value);
                    }
                }}
                clearable
                style={{ width: '300px' }}
            ></Select>
            <Button onClick={handleGrantPermissions}>Login</Button>
            <Button onClick={handleRemovePermissions}>Logout</Button>
        </>
    )
}

export default GrantPermissions;
