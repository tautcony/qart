<!DOCTYPE html>

<html dir="ltr" lang="{{.CurLang.Lang}}">
<head>
    <title>QArt Coder</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8" />
    <link rel="stylesheet" href="/static/css/bootstrap.min.css" />
    <style>
        .qr-container {
            margin-top: 3em;
            text-align: center;
        }

        .qr-container #op-qr-code, .qr-container .figure {
            width: 100%;
            max-width: 300px;
        }

        .parameter-container {
            margin-top: 1em;
        }

        #op-qr-code {
            cursor: pointer;
        }
    </style>
 </head>

<body>
    <div class="container">
        <div class="row">
            <div class="qr-container col-sm-6 col-md-4 col-lg-3">
                <figure class="figure">
                    <img id="op-qr-code"
                         class="figure-img img-fluid rounded"
                         src="/image/placeholder/400x400" alt="QR Code"
                         data-toggle="tooltip" data-placement="top" title='{{i18n .Lang "index.rotate"}}'
                    />
                </figure>
                <div class="row g-3">
                    <div class="col-6">
                        <button id="op-refresh" type="button" class="btn btn-primary">{{i18n .Lang "index.refresh"}}</button>
                    </div>
                    <div class="col-6">
                        <button id="op-share" type="button" class="btn btn-primary">{{i18n .Lang "index.share"}}</button>
                    </div>
                </div>
                <div class="row" style="padding-top: 1em">
                    <a href="https://research.swtch.com/qart" target="_blank" rel="noopener">{{i18n .Lang "index.how"}}</a>
                </div>
            </div>
            <div class="parameter-container col-sm-6 col-md-8 col-lg-9">
                <div class="row">
                    <h1 class="col">QArt Coder</h1>
                    <div class="dropdown col">
                        <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-expanded="false">
                            {{.CurLang.Name}}
                        </button>
                        <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton">
                            {{range $i, $v := .RestLangs}}
                                <li><a class="dropdown-item" href="/?lang={{$v.Lang}}">{{$v.Name}}</a></li>
                            {{end}}
                        </ul>
                    </div>
                </div>
                <form class="row g-3" id="parameter-form">
                    <div class="input-group mb-3 col-md-16">
                        <button id="op-upload" type="button" class="btn btn-primary file">
                            {{i18n .Lang "index.upload"}}
                        </button>
                        <input type="file" id="op-upload-input" accept="image/*" style="display:none" />
                        <input class="form-control" type="text" value="https://example.com" id="op-url" />
                    </div>
                    <div class="form-check form-switch col-md-6">
                        <input class="form-check-input" type="checkbox" value="" id="op-rand-control" />
                        <label class="form-check-label" for="op-rand-control">{{i18n .Lang "index.rand_control"}}</label>
                    </div>
                    <div class="form-check form-switch col-md-6">
                        <input class="form-check-input" type="checkbox" value="" id="op-only-data-bits" />
                        <label class="form-check-label" for="op-only-data-bits">{{i18n .Lang "index.only_data_bits"}}</label>
                    </div>
                    <div class="form-check form-switch col-md-6">
                        <input class="form-check-input" type="checkbox" value="" id="op-dither" />
                        <label class="form-check-label" for="op-dither">{{i18n .Lang "index.dither"}}</label>
                    </div>
                    <div class="form-check form-switch col-md-6">
                        <input class="form-check-input" type="checkbox" value="" id="op-save-control" />
                        <label class="form-check-label" for="op-save-control">{{i18n .Lang "index.save_control"}}</label>
                    </div>
                    <div class="col-6">
                        <label for="op-dx" class="form-label">X</label>
                        <input type="range" data-reverse="" class="form-range col-12" min="-50" max="50" step="1" value="-4" id="op-dx" />
                        <label for="op-dy" class="form-label">Y</label>
                        <input type="range" data-reverse="" class="form-range col-12" min="-50" max="50" step="1" value="-4" id="op-dy" />
                    </div>
                    <div class="col-6">
                        <label for="op-version" class="form-label">{{i18n .Lang "index.qr_version"}}</label>
                        <input type="range" class="form-range col-12" min="1" max="9" step="1" value="6" id="op-version" />
                        <label for="op-size" class="form-label">{{i18n .Lang "index.image_size"}}</label>
                        <input type="range" class="form-range col-12" min="-20" max="20" step="1" value="0" id="op-size" />
                    </div>
                </form>
            </div>
        </div>
    </div>

    <script src="/static/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/qart.js"></script>
    <!--script src="/static/js/reload.min.js"></script-->
    <a href="https://github.com/tautcony/qart" class="github-corner" aria-label="View source on GitHub"><svg width="80" height="80" viewBox="0 0 250 250" style="fill:#151513; color:#fff; position: absolute; top: 0; border: 0; right: 0;" aria-hidden="true"><path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path><path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2" fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path><path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z" fill="currentColor" class="octo-body"></path></svg></a><style>.github-corner:hover .octo-arm{animation:octocat-wave 560ms ease-in-out}@keyframes octocat-wave{0%,100%{transform:rotate(0)}20%,60%{transform:rotate(-25deg)}40%,80%{transform:rotate(10deg)}}@media (max-width:500px){.github-corner:hover .octo-arm{animation:none}.github-corner .octo-arm{animation:octocat-wave 560ms ease-in-out}}</style>
</body>
</html>
