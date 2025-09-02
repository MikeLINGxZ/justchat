const prefixErr = "@backend@err@";

export const packageErr = (errStr: string) => {
    return prefixErr+errStr;
}

export const unpackageErr = (errStr: string) => {
    return errStr.replace(prefixErr, "");
}

export const isError = (errStr: string) => {
    return errStr.indexOf(prefixErr) > -1;
}