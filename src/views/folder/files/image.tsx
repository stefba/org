import React from 'react';
import File from 'funcs/files';
import Info from './Info';
import AddInfo from './extra/AddInfo';

export default function Image({file, isSingle}: { file: File; isSingle?: boolean }) {
    if (isSingle) {
        return <img alt="" src={"/file" + file.path} />
    } else {
        return (
            <div className="list-media">
                <img alt="" src={"/file" + file.path} />
                <AddInfo mainFilePath={file.path}></AddInfo>
                <Info mainFilePath={file.path} />
            </div>
        )
    }
}


