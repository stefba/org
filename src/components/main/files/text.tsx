import React, { useState, useEffect, useRef } from 'react';
//import TextareaAutosize from 'react-textarea-autosize';
import TextareaAutosize from '@material-ui/core/TextareaAutosize';
import { Meta } from 'components/main/files/meta';
import { mainFuncsObj, modFuncsObj } from 'components/main/main';
import File from 'funcs/files';

type TextFieldProps = {
    file: File;
    mainFuncs?: mainFuncsObj;
    modFuncs: modFuncsObj;
    isSingle: boolean;
}

export default function TextField({file, mainFuncs, modFuncs, isSingle}: TextFieldProps) {
    const [body, setBody] = useState(file.body);

    const ref = useRef<HTMLTextAreaElement>(null!)

    useEffect(() => {

        setBody(file.body);

        if (isSingle && ref && ref.current) {
            ref.current.focus({preventScroll:true});
        }

    }, [file, isSingle]);

    function handleTyping(e: React.FormEvent<HTMLTextAreaElement>) {
        setBody(e.currentTarget.value);
    }

    function checkSubmit(e: React.KeyboardEvent<HTMLTextAreaElement>) {
        if (e.ctrlKey && e.key === "Enter" && mainFuncs !== undefined) {
            submit(e);
            mainFuncs.createNewFile();
        }
    }

    function submit(e: React.FormEvent<HTMLTextAreaElement>) {
        file.body = body;
        modFuncs.writeFile(file);
    }

    return (
        <div className={"text" + (isNoSort(file.name) ? " no-sort" : "")}>
        { !isSingle &&
            <Meta file={file} modFuncs={modFuncs}/>
        }
        <TextareaAutosize value={body}
        ref={ref}
        rowsMin={!isSingle ? 1 : fullScreenRows()}
        onKeyPress={checkSubmit}
        onChange={handleTyping}
        onBlur={submit} />
        </div>
    )
}

function isNoSort(name: string): boolean {
    switch (name) {
        case "info":
        case ".sort":
            return true;
        default:
            return false;
    }
}

function fullScreenRows() {
    return Math.round(window.screen.height/(2.25*16)) - 8
}

