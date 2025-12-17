import React from 'react'
import { HashRouter } from 'react-router-dom'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import './i18n/config'; // Initialize i18n

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <HashRouter>
            <App/>
        </HashRouter>
    </React.StrictMode>
)
