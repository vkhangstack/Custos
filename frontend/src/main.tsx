import React from 'react'
import { HashRouter } from 'react-router-dom'
import { createRoot } from 'react-dom/client'
import './style.css'
import App from './App'
import './i18n/config'; // Initialize i18n
import { ThemeProvider } from './context/ThemeContext'
import { ToastProvider } from './context/ToastContext'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
            <HashRouter>
                <ToastProvider>
                    <App />
                </ToastProvider>
            </HashRouter>
        </ThemeProvider>
    </React.StrictMode>
)
