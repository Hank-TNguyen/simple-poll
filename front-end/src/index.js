import React from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App';
// Remove the below line if you don't have reportWebVitals.js
import reportWebVitals from './reportWebVitals';

const container = document.getElementById('root');
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// If removed, also remove or comment out the reportWebVitals() call.
reportWebVitals(console.log);