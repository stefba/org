import React from 'react';
import 'css/main.scss';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import TargetsProvider from './context/targets';
import ErrProvider from './context/err';
import Write from 'views/Write';
import Search from 'views/Search';
import Today from 'views/Today';
import Topics, { Topic } from 'views/Topics';
import Loader from 'Loader';

export default function App() {
    return (
    <TargetsProvider>
    <ErrProvider>
        <Router>
            <Routes>
                <Route path="/write" element={<Write />} />
                <Route path="/today" element={<Today />} />
                <Route path="/search/*" element={<Search />} />
                <Route path="/topics" element={<Topics />} />
                <Route path="/topics/*" element={<Topic />} />
                <Route path="/*" element={<Loader />} />
            </Routes>
        </Router>
    </ErrProvider>
    </TargetsProvider>
    )
}