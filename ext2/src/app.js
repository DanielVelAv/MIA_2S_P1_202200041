import './App.css';
import { PPrincipal } from './components/PPrincipal';
import React, { useState, useEffect } from 'react';
import axios from 'axios';

function App() {

    const [message, setMessage] = useState('');

    useEffect(() => {
        axios.get('http://localhost:8080/api/hello')
        .then(response => {
            setMessage(response.data.message);
        });
    }, []);

    return (
        <div className="App">
        <header className='App-header'>
        <PPrincipal />
        </header>
        </div>
    
    );
}

export default App;