import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Log from './components/Log';
import Config from './components/Config';
import Network from './components/Network';

import { useTranslation } from 'react-i18next';

function Home() {
    const { t } = useTranslation();
    const [resultText, setResultText] = useState(t('home.welcome'));
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>{t('home.greetButton')}</button>
            </div>
        </div>
    )
}

function App() {
    return (
        <Routes>
            <Route path="/" element={<Layout />}>
                <Route index element={<Home />} />
                <Route path="log" element={<Log />} />
                <Route path="config" element={<Config />} />
                <Route path="network" element={<Network />} />
            </Route>
        </Routes>
    )
}

export default App
