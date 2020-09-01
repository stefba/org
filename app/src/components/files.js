import React, { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { basename } from 'path';
import TextareaAutosize from 'react-textarea-autosize';

const File = ({ file }) => {
  switch (file.type) {
    case "dir":
      return Dir({file})
    case "text":
      return Text({file})
    default:
      return Dir({file})
  }
}

const Text = ({file}) => {
  const [body, setBody] = useState("");

  useEffect(() => {
    fetch("/api/view" + file.path).then(
      resp => resp.text().then(
        textContent => setBody(textContent)
      )
    );

  }, [file]);

  return (
    <>
      <Dir file={file} />
      <TextareaAutosize value={body} />
    </>
  )
}

const Dir = ({file}) => {
  return (
    <>
      <div>- <Link to={file.path}>{basename(file.path)}</Link> ({file.type})</div>
    </>
  )
}

const Files = () => {
  const [files, setFiles] = useState([]);

  let path = useLocation().pathname;
  useEffect(() => {
    fetch("/api/view" + path).then(
      resp => resp.json().then(
        files => setFiles(files)
      ));

  }, [path]);

  return (
    files.map((file, i) => (
      <File key={i} file={file} />
    ))
  );
}

export default Files;


