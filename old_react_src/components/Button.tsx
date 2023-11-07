interface ButtonProps {
    children: string;
    colour?: 'primary' | 'secondary' | 'danger';
    onClick: () => void;
}

const Button = ({ children, onClick, colour = 'primary' }: ButtonProps) => {
    return (
        <>
            <h1>Buttons</h1>
            <button type="button" className={'btn btn-' + colour} onClick={onClick}>{children}</button>
        </>
    )
}

export default Button