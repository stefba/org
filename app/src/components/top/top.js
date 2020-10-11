import React, { useState, useEffect, useCallback } from 'react';
import { Link, useHistory } from 'react-router-dom';
import ThemeIcon from '@material-ui/icons/WbSunnySharp';
import TargetIcon from '@material-ui/icons/VerticalAlignBottom';
import { basename } from 'path';
import CrumbNav from './crumbs';
import { Del } from '../meta';
import * as p from '../../funcs/paths';
import * as targets from '../../funcs/targets';

const DirName = (path) => {
  const name = basename(path);
  if (name === "") {
    return "org"
  }
  return name
}

const PageTitle = path => {
  if (path === "/") {
    return "ORG"
  }
  return DirName(path) + " - ORG";
}

const makeNeighborList = files => {
  let nu = [];
  for (const f of files) {
    if (f.type === "dir") {
      nu.push(f.path);
    }
  }
  return nu;
}

const readStateBool = key => {
  const str = localStorage.getItem(key);
  if (str == null) {
    return false;
  }
  return str === "true";
}

const Top = ({view}) => {

  /* theme */

  const [darkTheme, setDarkTheme] = useState(readStateBool("dark-theme"));

  useEffect(() => {
    darkTheme
      ? document.body.dataset["theme"] = "dark"
      : document.body.dataset["theme"] = ""
  })

  function toggleTheme() {
    setDarkTheme(!darkTheme);
    localStorage.setItem("dark-theme", !darkTheme);
  }

  /* links */

  const [links, setLinks] = useState([]);

  useEffect(() => {
    fetch("/api/links").then(
      resp => resp.json().then(
        links => setLinks(links)
      )
    )
  }, [])

  /* targets */

  const [targetList, setTargetList] = useState([]);
  const [activeTarget, setActiveTarget] = useState("");

  const loadTargets = () => {
    setActiveTarget(targets.getActive());
    setTargetList(targets.getList())
  }

  const listenForTargets = useCallback(evt => {
    loadTargets();
  }, []);

  useEffect(() => {
    window.addEventListener('storage', listenForTargets);

    return () => {
      window.removeEventListener('storage', listenForTargets);
    };
  }, [listenForTargets]);

  useEffect(() => {
    loadTargets();
  }, [])

  /* targets functions */

  const setActive = path => {
    targets.setActive(path);
    loadTargets();
  }

  const removeTarget = path => {
    targets.removeTarget(path);
    loadTargets();
  }

  const setThisActive = () => {
    targets.setActive(view.file.path);
    loadTargets();
  }

  const TargetButton = ({clickFn}) => {
    return <button onClick={clickFn}><TargetIcon /></button>
  }


  /* path */

  const [path, setPath] = useState(view.file.path);

  useEffect(() => {
    setPath(view.file.path);
  }, [view]);

  document.title = PageTitle(path);

  const history = useHistory();

  async function renameFile(newPath) {
    try {
      const resp = await fetch("/api" + path, {
        method: "PUT",
        body: newPath
      });
      if (!resp.ok) {
        alert( "Rename failed: " + path + "\nreason: " +resp.statusText);
        return;
      }
      history.push(newPath)
    } catch(err) {
      console.log(err);
    }
  }

  async function delFile(path) {
    try {
      const resp = await fetch("/api" + path, { method: "DELETE" });
      if (!resp.ok) {
        alert( "Delete failed: " + path + "\nreason: " +resp.statusText);
        return;
      }
      history.push(view.parent);
    } catch(err) {
      console.log(err);
    }
  }

  return (
    <>
      <nav id="links">
        <span className="links__top">
          <LinkList links={links} />
        </span>
        <span id="targets" className="right">
          <TargetList
            links={targetList}
            page={path}
            activeTarget={activeTarget}
            setActiveFn={setActive}
            removeFn={removeTarget} />
        </span>
      </nav>

      <nav id="bar">
        <CrumbNav
          neighbors={makeNeighborList(view.neighbors)}
          switchLink={view.switch}
          path={path} />
        <span className="right">
          <TargetButton clickFn={setThisActive} />
          <button onClick={toggleTheme} ><ThemeIcon /></button>
          <Del file={view.file} delFile={delFile} />
        </span>
      </nav>

      <h1 className="name">
        <Link className="parent" to={view.parent}>^</Link>
        <RenameInput path={path} renameFile={renameFile} />
      </h1>
    </>
  )
}

const RenameInput = ({path, renameFile}) => {
  const [name, setName] = useState(DirName(path));

  useEffect(() => {
    setName(DirName(path));
  }, [path]);

  function handleTyping(evt) {
    setName(evt.target.value);
  }

  function submit() {
    const old = DirName(path);
    if (old === name) {
      return;
    }
    const dir = p.Dir(path);
    renameFile(dir + (dir === "/" ? "" : "/") + name);
  }

  return (
    <input type="text" value={name}
      disabled={name === "org" ? "disabled" : ""}
      onChange={handleTyping} onBlur={submit} />
  )
}

const TargetList = ({activeTarget, page, links, removeFn, setActiveFn}) => {
  const onClick = evt => {
    if (evt.shiftKey) {
      evt.preventDefault();
      setActiveFn(evt.target.pathname);
    }
  }
  const onRightClick = evt => {
    if (evt.shiftKey) {
      evt.preventDefault();
      removeFn(evt.target.pathname);
    }
  }
  return (
    links.map((l, i) => {
      let className = p.Section(l)
      if (activeTarget === l) {
        className += " active"
      }
      return (
        <Link key={i} to={l}
        className={className}
        onClick={onClick} onContextMenu={onRightClick}>
          {p.ExtendedBase(l)}
        </Link>
      )
    })
  )
}

const LinkList = ({links, active}) => {
  if (!links || links.length === 0) {
    return null;
  }
  return (
    links.map((l, i) => (
      <Link key={i} to={l} className={active === l ? "active" : ""}>{p.Base(l)}</Link>
    ))
  )
}

export default Top;