
import React from 'react';

interface IProps {
    message: string
}

const SimpleComponent: React.FunctionComponent<IProps> = ({ message }) => {
    return (
        <div>{message}</div>
    )
}

export default SimpleComponent



let obj = { hej: 'meddig' };

let { hej } = obj;

console.log(hej);

function asd(optionObject) {
    let host = optionObject.host;
    let port = optionObject.port || 3306;
    let username = optionObject.username;
}

asd({
    host: '1.2.3.4',
    port: 3306,
    username: 'jrl'
})

// BETTER

const asdd = ({ host, port = 3306, username }) => {
    // thats it
}