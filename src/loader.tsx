import React, { useState, useEffect, useCallback, useContext } from 'react';
import { useLocation } from 'react-router-dom';
import { TargetsContext } from './context/targets';
import Nav from 'parts/Nav';
import { isDir, pageTitle } from 'funcs/paths';
import { setFavicon, blinkFavicon } from 'funcs/favicon';
import { isActiveTarget } from 'funcs/targets';
import { dirPath } from 'funcs/paths';
import useView, { viewObject } from 'state';
import View from 'views/View';

export default function Loader() {
    const { targets } = useContext(TargetsContext);
    const path = useLocation().pathname;

    const { view, setView } = useView()

    const [status, setStatus] = useState("");

    async function loadView(path: string) {
        const resp = await fetch("/api/view" + path);
        if (!resp.ok) {
            setStatus(resp.status + " - " + resp.statusText)
        }
        const newView = await resp.json();
        setView(newView);
    };

    // only load when a new *dir* is requested
    useEffect(() => {
        document.title = pageTitle(path);
        setFavicon(path);

        if (shouldLoad(path, view)) {
            loadView(path);
        }
    }, [path, view]);

    function shouldLoad(path: string, view: viewObject): boolean {
        // load if is new dir
        if (isDir(path) && dirPath(path) !== view.path) {
            return true
        }
        // load if is file but no dir loaded
        if (!isDir(path) && view.path === "") {
            return true
        }
        return false
    }

    // listen if another tab sends files to this tab.
    const listenForWrite = useCallback(() => {
        if (isActiveTarget(targets, path)) {
            loadView(path);
            blinkFavicon(path);
        }
    }, [targets, path]);

    useEffect(() => {
        window.addEventListener('storage', listenForWrite);

        return () => {
            window.removeEventListener('storage', listenForWrite);
        };
    }, [listenForWrite]);

    if (status !== "") {
        return (
            <>
                <Nav path={path} />
                <br />
                <code>{status}</code>
            </>
        );
    }

    return (
        <>
            <Nav path={path} />
            <View path={path} files={view.dir.files} sorted={view.dir.sorted} />
        </>
    );
}
