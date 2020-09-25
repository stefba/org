import React from 'react';
import { Link, useLocation, useHistory } from 'react-router-dom';
import { basename } from 'path';
import Breadcrumbs from './breadcrumbs';
import { Del } from './info';

const DirName = () => {
  const name = basename(useLocation().pathname);
  if (name === "") {
    return "org"
  }
  return name
}

const Root = () => {
  return (
    <Link to="/">org</Link>
  )
}

const Top = ({view}) => {

  /*
  const rename = newName => {
  }
  */

  const history = useHistory();

  const del = path => {
    fetch("/api" + path, {
      method: "DELETE"
    }).then(resp => {
      if (!resp.ok) {
        resp.text().then(
          msg => {
            alert(msg)
            return;
            //throw Error(msg);
          }
        );
        return;
      };
      history.push(view.parent);
    }).catch(err => {
      alert(err);
      console.log(err);
    })
  }


  return (
    <>
      <nav>
        <Root />
        <Breadcrumbs />
        <span className="right">
          <Del file={view.file} delFn={del} />
        </span>
      </nav>
      <h1>
        <Link to={view.parent}>^</Link>
        <DirName />
      </h1>
    </>
  )
}

export default Top;