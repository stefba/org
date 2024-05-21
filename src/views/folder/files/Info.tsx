import React from 'react';

import TextField from './TextView';
import { newInfoFile } from 'funcs/files';

export default function Info({ mainFilePath }: { mainFilePath: string }) {
    const info = newInfoFile(mainFilePath);
    return (
        <TextField file={info} isSingle={false} />
    )
}
