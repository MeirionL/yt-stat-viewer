import { Button, TextInput } from '@mantine/core';
import { ChangeEvent, useState } from 'react';

function GrantPermissions() {
    const [userId, setUserId] = useState<string>('');

    const handleGrantPermissions = () => {
        window.location.href = 'http://localhost:8080/auth';
    };

    const handleRemovePermissions = () => {
        if (userId) {
            window.location.href = `http://localhost:8080/logout/${userId}`;
        } else {
            console.error('Please enter a user ID');
        }
    };

    const handleUserIdChange = (event: ChangeEvent<HTMLInputElement>) => {
        setUserId(event.target.value);
    };

    return (
        <>
            <Button onClick={handleGrantPermissions}>Login</Button>
            <TextInput
                label="Enter ID"
                description="Enter ID for account you want to log out"
                value={userId}
                onChange={(handleUserIdChange)}
                size='sm'
                style={{ marginBottom: '8px' }}
            />
            <Button onClick={handleRemovePermissions}>Logout</Button>
        </>
    );
}

export default GrantPermissions;
