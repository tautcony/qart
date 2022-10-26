function timeoutPromise(ms, promise) {
    return new Promise((resolve, reject) => {
        const timeoutId = setTimeout(() => {
            reject(new Error("promise timeout"))
        }, ms);
        promise.then(
            (res) => {
                clearTimeout(timeoutId);
                resolve(res);
            },
            (err) => {
                clearTimeout(timeoutId);
                reject(err);
            }
        );
    })
}

function errorHandler(response) {
    if (response && !response.success) {
        alert(response.message);
    }
    return response
}

function request(input, init) {
    return fetch(input, init).then(response => response.json()).then(errorHandler);
}

const Api = {
    render: operation => {
        return request('/v1/render', {
            method: 'POST',
            body: JSON.stringify(operation)
        }).then(response => {
            const element = document.getElementById('op-qr-code');
            if (element) {
                element.src = response.data.image;
            }
        });
    },
    config: () => request('/v1/render/config'),
    upload: file => {
        let data = new FormData();
        data.append('image', file);

        return request('/v1/render/upload', {
            method: 'POST',
            body: data
        });
    },
    share: operation => request('/v1/share', {
        method: 'POST',
        body: JSON.stringify({
            image: operation.image
        })
    })
}

const Element = {
    img: {
        op_qr_code: document.getElementById('op-qr-code'),
    },
    button: {
        op_refresh: document.getElementById('op-refresh'),
        op_share: document.getElementById('op-share'),
        op_upload: document.getElementById('op-upload'),
    },
    input: {
        op_upload_input: document.getElementById('op-upload-input'),
        op_url: document.getElementById('op-url'),
        op_rand_control: document.getElementById('op-rand-control'),
        op_only_data_bits: document.getElementById('op-only-data-bits'),
        op_dither: document.getElementById('op-dither'),
        op_save_control: document.getElementById('op-save-control'),
        op_dx: document.getElementById('op-dx'),
        op_dy: document.getElementById('op-dy'),
        op_version: document.getElementById('op-version'),
        op_size: document.getElementById('op-size'),
    }
}

function updateOperation(element, obj) {
    if (!obj || !element.id || !element.id.startsWith('op-')) {
        return;
    }
    let id = element.id.replace('op-', '').replace(/-/g, '').toLocaleLowerCase();
    if (!(id in obj)) {
        return;
    }
    switch (element.type) {
        case 'text': {
            obj[id] = element.value;
            break;
        }
        case 'checkbox': {
            obj[id] = element.checked;
            break;
        }
        case 'range': {
            let value = parseInt(element.value, 10);
            obj[id] = 'reverse' in element.dataset ? -value : value;
        }
    }
    Api.render(obj);
}

function applyOperation(operation) {
    Element.input.op_url.value = operation.url;
    Element.input.op_rand_control.checked = operation.randcontrol;
    Element.input.op_only_data_bits.checked = operation.onlydatabits;
    Element.input.op_dither.checked = operation.dither;
    Element.input.op_save_control.checked = operation.savecontrol;
    Element.input.op_dx.value = -operation.dx;
    Element.input.op_dy.value = -operation.dy;
    Element.input.op_version.value = operation.version;
    Element.input.op_size.value = operation.size;
}

(() => {
    let operation = {
        image: "default",
        dx: 4,
        dy: 4,
        size: 0,
        url: "https://example.com",
        version: 6,
        mask: 2,
        randcontrol: false,
        dither: false,
        onlydatabits: false,
        savecontrol: false,
        seed: "",
        scale: 4,
        rotation: 0
    };
    timeoutPromise(500, Api.config().then(response => {
        operation = response.data;
    })).catch(ignore => {}).then(() => {
        applyOperation(operation);
    });
    Object.values(Element.input).forEach(element => {
        let handler = event => {
            updateOperation(event.target, operation);
        };
        if (element.type === 'text') {
            element.addEventListener('input', handler, false);
        } else {
            element.addEventListener('change', handler);
        }
    });
    Element.input.op_upload_input.addEventListener('input', event => {
        let files = event.target.files;
        if (files && files.length > 0) {
            Api.upload(files[0]).then(function (response) {
                operation.image = response.data.id;
                Api.render(operation);
            });
        }
    }, false);
    Element.img.op_qr_code.addEventListener('click', event => {
        operation.rotation = (operation.rotation + 1) % 4;
        Api.render(operation);
    }, false);
    Element.button.op_refresh.addEventListener('click', event => Api.render(operation), false);
    Element.button.op_upload.addEventListener('click', event => Element.input.op_upload_input.click(), false);
    Element.button.op_share.addEventListener('click', event => {
        Api.share(operation).then(response => {
            window.open(`/share/${response.data.id}`);
        });
    });
})();
