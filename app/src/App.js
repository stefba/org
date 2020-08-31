import React, {useState, useEffect} from 'react';
import './App.css';

const File = ({ file }) => {
  console.log("x");
  return (
    <div>{file.path}</div>
  )
}

const App = () => {
  const [files, setFiles] = useState([]);

  useEffect(() => {

    fetch("/api/view/").then(
      resp => resp.json().then(
        files => setFiles(files)
      ));

  }, []);

  return (
    files.map((file, i) => (
      <File key={i} file={file} />
    ))
  );
}

export default App;
