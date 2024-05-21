import React from 'react';
import File from 'funcs/files';
import AddInfo from 'views/folder/files/extra/AddInfo'

type VideoProps = {
    file: File;
    isSingle?: boolean;
}

export default function Video({file, isSingle}: VideoProps) {
    if (isSingle) {
        return (
            <video controls src={"/file" + file.path}></video>
        )
    } else {
        return (
            <div className="list-media">
                <video controls src={"/file" + file.path}></video>
                <AddInfo mainFilePath={file.path}></AddInfo>
            </div>
        )
    }
}

