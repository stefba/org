import React, { useEffect, useContext, useState } from 'react';
import { Link } from 'react-router-dom';
import ReverseIcon from '@material-ui/icons/SwapVert';
import { ReactSortable } from 'react-sortablejs';
import { basename } from 'path';
import { Info } from 'components/meta';
import Text from 'components/types/text';
import Image from 'components/types/image';
import { ModFuncs } from 'components/main/view';
import File from 'funcs/file';
import { setActiveTarget } from 'funcs/targets';
import { TargetsContext } from 'context/targets';

type FileSwitchProps = {
    file: File;
    modFuncs: ModFuncs;
    isSingle: boolean;
}

function FileSwitch({file, modFuncs, isSingle}: FileSwitchProps) {
    /*
  if (file.name === ".sort") {
    return <div className="no-sort"><Info file={file} modFuncs={modFuncs} /></div>
  }
     */
    switch (file.type) {
        case "text":
            return <Text file={file} modFuncs={modFuncs} isSingle={isSingle}/>
        case "image":
            return <Image file={file} modFuncs={modFuncs} />
        default:
            return <Info file={file} modFuncs={modFuncs} />
    }
}

type FileListProps = {
    files: File[];
    saveSort: (part: File[], type: string) => void;
    modFuncs: ModFuncs;
}

function FileList({files, saveSort, modFuncs}: FileListProps) {
    const [state, setState] = useState(files);

    useEffect(() => {
        setState(files);
    }, [files])

    const reverseFiles = () => {
        const reverse = preSort(state.slice().reverse());
        saveSort(reverse, "files");
    }

    const callOnEnd = () => {
        saveSort(state, "files");
    };

    if (!files || files.length === 0) {
        return null
    }
    // filter=".no-sort"
    return (
        <>
            <span className="right">
                <button onClick={reverseFiles}><ReverseIcon /></button>
            </span>
            <ReactSortable 
                handle=".info__drag" 
                onEnd={callOnEnd}
                animation={200} list={state} setList={setState}>

                    { state.map((file, i) => (
                        <FileSwitch key={file.id} file={file} modFuncs={modFuncs} isSingle={false} />
                    ))}

            </ReactSortable>
        </>
    );
}

function Dir({dir}: {dir: File;}) {
    let { targets, saveTargets } = useContext(TargetsContext);

    function setTarget(e: React.MouseEvent<HTMLAnchorElement>) {
        if (e.shiftKey) {
            console.log("try to set something");
            e.preventDefault();
            saveTargets(setActiveTarget(targets, e.currentTarget.pathname));
        }
    }
    return (
        <Link to={dir.path} onClick={setTarget}>{basename(dir.path)}</Link>
    )
}

type DirListProps = {
    dirs: File[];
    saveSort: (part: File[], type: string) => void;
}

function DirList({dirs, saveSort}: DirListProps) {
    const [state, setState] = useState(dirs);

    useEffect(() => {
        setState(dirs);
    }, [dirs])

    const callOnEnd = () => {
        saveSort(state, "dirs");
    };

    return (
        <ReactSortable className="dirs__list" onEnd={callOnEnd}
        animation={200} list={state} setList={setState}>
        {state.map((dir) => (
            <Dir key={dir.id} dir={dir} />
        ))}
        </ReactSortable>
    )
}

export { DirList, FileList, FileSwitch };

function preSort(files: File[]): File[] {
    let info = [];
    let sort = [];
    let nu   = [];

    for (const f of files) {
        if (f.name === "info") {
            info.push(f)
            continue
        }
        if (f.name === ".sort") {
            sort.push(f)
            continue
        }
        nu.push(f)
    }

    return info.concat(nu).concat(sort);
}
