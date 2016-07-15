getData = function (func) {
    var args = [].slice.call(arguments);
    args.shift();
    return new Promise((resolve, reject) => {
        args.push((error, data) => {
            if (error) {
                reject(error);
            }
            resolve(data);
        });
        func.apply(null, args);
    });
};