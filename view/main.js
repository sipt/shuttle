(window["webpackJsonp"] = window["webpackJsonp"] || []).push([["main"],{

/***/ "./src/$$_lazy_route_resource lazy recursive":
/*!**********************************************************!*\
  !*** ./src/$$_lazy_route_resource lazy namespace object ***!
  \**********************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

function webpackEmptyAsyncContext(req) {
	// Here Promise.resolve().then() is used instead of new Promise() to prevent
	// uncaught exception popping up in devtools
	return Promise.resolve().then(function() {
		var e = new Error("Cannot find module '" + req + "'");
		e.code = 'MODULE_NOT_FOUND';
		throw e;
	});
}
webpackEmptyAsyncContext.keys = function() { return []; };
webpackEmptyAsyncContext.resolve = webpackEmptyAsyncContext;
module.exports = webpackEmptyAsyncContext;
webpackEmptyAsyncContext.id = "./src/$$_lazy_route_resource lazy recursive";

/***/ }),

/***/ "./src/app/app-routing.module.ts":
/*!***************************************!*\
  !*** ./src/app/app-routing.module.ts ***!
  \***************************************/
/*! exports provided: AppRoutingModule */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "AppRoutingModule", function() { return AppRoutingModule; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/router */ "./node_modules/@angular/router/fesm5/router.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};


var routes = [];
var AppRoutingModule = /** @class */ (function () {
    function AppRoutingModule() {
    }
    AppRoutingModule = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["NgModule"])({
            imports: [_angular_router__WEBPACK_IMPORTED_MODULE_1__["RouterModule"].forRoot(routes)],
            exports: [_angular_router__WEBPACK_IMPORTED_MODULE_1__["RouterModule"]]
        })
    ], AppRoutingModule);
    return AppRoutingModule;
}());



/***/ }),

/***/ "./src/app/app.component.css":
/*!***********************************!*\
  !*** ./src/app/app.component.css ***!
  \***********************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ".st-layout {\n    height: 100vh;\n}\n.st-white-bgc {\n    background: #fff;\n}\n.st-footer{\n    position:fixed;\n}"

/***/ }),

/***/ "./src/app/app.component.html":
/*!************************************!*\
  !*** ./src/app/app.component.html ***!
  \************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<!-- NG-ZORRO -->\n<nz-layout class=\"st-layout\">\n  <nz-sider>\n      <app-menus></app-menus>\n  </nz-sider>\n  <nz-layout>\n    <nz-content class=\"st-white-bgc\">\n      <router-outlet></router-outlet>\n    </nz-content>\n  </nz-layout>\n</nz-layout>"

/***/ }),

/***/ "./src/app/app.component.ts":
/*!**********************************!*\
  !*** ./src/app/app.component.ts ***!
  \**********************************/
/*! exports provided: AppComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "AppComponent", function() { return AppComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};

var AppComponent = /** @class */ (function () {
    function AppComponent() {
        this.title = 'shuttle';
    }
    AppComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-root',
            template: __webpack_require__(/*! ./app.component.html */ "./src/app/app.component.html"),
            styles: [__webpack_require__(/*! ./app.component.css */ "./src/app/app.component.css")]
        })
    ], AppComponent);
    return AppComponent;
}());



/***/ }),

/***/ "./src/app/app.module.ts":
/*!*******************************!*\
  !*** ./src/app/app.module.ts ***!
  \*******************************/
/*! exports provided: AppModule */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "AppModule", function() { return AppModule; });
/* harmony import */ var _angular_platform_browser__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/platform-browser */ "./node_modules/@angular/platform-browser/fesm5/platform-browser.js");
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _app_routing_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./app-routing.module */ "./src/app/app-routing.module.ts");
/* harmony import */ var _app_component__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./app.component */ "./src/app/app.component.ts");
/* harmony import */ var _components_records_records_component__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./components/records/records.component */ "./src/app/components/records/records.component.ts");
/* harmony import */ var _components_dns_cache_dns_cache_component__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! ./components/dns-cache/dns-cache.component */ "./src/app/components/dns-cache/dns-cache.component.ts");
/* harmony import */ var _angular_platform_browser_animations__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! @angular/platform-browser/animations */ "./node_modules/@angular/platform-browser/fesm5/animations.js");
/* harmony import */ var _angular_forms__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(/*! @angular/forms */ "./node_modules/@angular/forms/fesm5/forms.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var ng_zorro_antd__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(/*! ng-zorro-antd */ "./node_modules/ng-zorro-antd/esm5/antd.js");
/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_10__ = __webpack_require__(/*! @angular/common */ "./node_modules/@angular/common/fesm5/common.js");
/* harmony import */ var _angular_common_locales_zh__WEBPACK_IMPORTED_MODULE_11__ = __webpack_require__(/*! @angular/common/locales/zh */ "./node_modules/@angular/common/locales/zh.js");
/* harmony import */ var _angular_common_locales_zh__WEBPACK_IMPORTED_MODULE_11___default = /*#__PURE__*/__webpack_require__.n(_angular_common_locales_zh__WEBPACK_IMPORTED_MODULE_11__);
/* harmony import */ var _components_menus_menus_component__WEBPACK_IMPORTED_MODULE_12__ = __webpack_require__(/*! ./components/menus/menus.component */ "./src/app/components/menus/menus.component.ts");
/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_13__ = __webpack_require__(/*! @angular/router */ "./node_modules/@angular/router/fesm5/router.js");
/* harmony import */ var _components_general_general_component__WEBPACK_IMPORTED_MODULE_14__ = __webpack_require__(/*! ./components/general/general.component */ "./src/app/components/general/general.component.ts");
/* harmony import */ var _components_server_server_component__WEBPACK_IMPORTED_MODULE_15__ = __webpack_require__(/*! ./components/server/server.component */ "./src/app/components/server/server.component.ts");
/* harmony import */ var _pipes_ips_format_pipe__WEBPACK_IMPORTED_MODULE_16__ = __webpack_require__(/*! ./pipes/ips-format.pipe */ "./src/app/pipes/ips-format.pipe.ts");
/* harmony import */ var _pipes_nl2br_pipe__WEBPACK_IMPORTED_MODULE_17__ = __webpack_require__(/*! ./pipes/nl2br.pipe */ "./src/app/pipes/nl2br.pipe.ts");
/* harmony import */ var _pipes_html2text_pipe__WEBPACK_IMPORTED_MODULE_18__ = __webpack_require__(/*! ./pipes/html2text.pipe */ "./src/app/pipes/html2text.pipe.ts");
/* harmony import */ var _components_mitm_mitm_component__WEBPACK_IMPORTED_MODULE_19__ = __webpack_require__(/*! ./components/mitm/mitm.component */ "./src/app/components/mitm/mitm.component.ts");
/* harmony import */ var _pipes_capacity_pipe__WEBPACK_IMPORTED_MODULE_20__ = __webpack_require__(/*! ./pipes/capacity.pipe */ "./src/app/pipes/capacity.pipe.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};





















Object(_angular_common__WEBPACK_IMPORTED_MODULE_10__["registerLocaleData"])(_angular_common_locales_zh__WEBPACK_IMPORTED_MODULE_11___default.a);
var appRoutes = [
    {
        path: '',
        redirectTo: '/records',
        pathMatch: 'full'
    },
    {
        path: 'records',
        component: _components_records_records_component__WEBPACK_IMPORTED_MODULE_4__["RecordsComponent"]
    },
    {
        path: 'dns-cache',
        component: _components_dns_cache_dns_cache_component__WEBPACK_IMPORTED_MODULE_5__["DnsCacheComponent"]
    },
    {
        path: 'servers',
        component: _components_server_server_component__WEBPACK_IMPORTED_MODULE_15__["ServerComponent"]
    },
    {
        path: 'mitm',
        component: _components_mitm_mitm_component__WEBPACK_IMPORTED_MODULE_19__["MitmComponent"]
    }
];
var AppModule = /** @class */ (function () {
    function AppModule() {
    }
    AppModule = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_1__["NgModule"])({
            declarations: [
                _app_component__WEBPACK_IMPORTED_MODULE_3__["AppComponent"],
                _components_records_records_component__WEBPACK_IMPORTED_MODULE_4__["RecordsComponent"],
                _components_dns_cache_dns_cache_component__WEBPACK_IMPORTED_MODULE_5__["DnsCacheComponent"],
                _components_menus_menus_component__WEBPACK_IMPORTED_MODULE_12__["MenusComponent"],
                _components_general_general_component__WEBPACK_IMPORTED_MODULE_14__["GeneralComponent"],
                _components_server_server_component__WEBPACK_IMPORTED_MODULE_15__["ServerComponent"],
                _pipes_ips_format_pipe__WEBPACK_IMPORTED_MODULE_16__["IpsFormatPipe"],
                _pipes_nl2br_pipe__WEBPACK_IMPORTED_MODULE_17__["Nl2BrPipe"],
                _pipes_html2text_pipe__WEBPACK_IMPORTED_MODULE_18__["Html2textPipe"],
                _components_mitm_mitm_component__WEBPACK_IMPORTED_MODULE_19__["MitmComponent"],
                _pipes_capacity_pipe__WEBPACK_IMPORTED_MODULE_20__["CapacityPipe"]
            ],
            imports: [
                _angular_platform_browser__WEBPACK_IMPORTED_MODULE_0__["BrowserModule"],
                _app_routing_module__WEBPACK_IMPORTED_MODULE_2__["AppRoutingModule"],
                _angular_platform_browser_animations__WEBPACK_IMPORTED_MODULE_6__["BrowserAnimationsModule"],
                _angular_forms__WEBPACK_IMPORTED_MODULE_7__["FormsModule"],
                _angular_common_http__WEBPACK_IMPORTED_MODULE_8__["HttpClientModule"],
                ng_zorro_antd__WEBPACK_IMPORTED_MODULE_9__["NgZorroAntdModule"],
                _angular_router__WEBPACK_IMPORTED_MODULE_13__["RouterModule"].forRoot(appRoutes)
            ],
            providers: [{ provide: ng_zorro_antd__WEBPACK_IMPORTED_MODULE_9__["NZ_I18N"], useValue: ng_zorro_antd__WEBPACK_IMPORTED_MODULE_9__["zh_CN"] }],
            bootstrap: [_app_component__WEBPACK_IMPORTED_MODULE_3__["AppComponent"]]
        })
    ], AppModule);
    return AppModule;
}());



/***/ }),

/***/ "./src/app/components/dns-cache/dns-cache.component.css":
/*!**************************************************************!*\
  !*** ./src/app/components/dns-cache/dns-cache.component.css ***!
  \**************************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ".st-ext-title {\n    background: #fafafa;\n    width: 100%;\n    height: 40px;\n}"

/***/ }),

/***/ "./src/app/components/dns-cache/dns-cache.component.html":
/*!***************************************************************!*\
  !*** ./src/app/components/dns-cache/dns-cache.component.html ***!
  \***************************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<div style=\"height: calc(100vh - 40px);\">\n  <nz-table #list [nzData]=\"cacheList\"\n  [nzShowPagination]=\"false\"\n  [nzFrontPagination]=\"false\"\n  nzSize=\"small\"\n  [nzScroll]=\"{y: 'calc(100vh - 80px)'}\">\n    <thead>\n      <tr>\n        <th nzWidth=\"100px\">Type</th>\n        <th nzWidth=\"300px\">DNSs</th>\n        <th nzWidth=\"200px\">Domain</th>\n        <th nzWidth=\"100px\">Country</th>\n        <th >IPs</th>\n      </tr>\n    </thead>\n    <tbody class=\"st-tbody\">\n      <tr *ngFor=\"let cache of list.data\">\n        <td>{{cache.Type}}</td>\n        <td>{{cache.DNSs}}</td>\n        <td>{{cache.Domain}}</td>\n        <td>{{cache.Country}}</td>\n        <td>\n          <nz-tag *ngFor=\"let ip of cache.IPs\">{{ip}}</nz-tag>\n        </td>\n      </tr>\n    </tbody>\n  </nz-table>\n</div>\n<div class=\"st-ext-title\">\n  <div>\n    <button style=\"margin: 4px;\" nz-button (click)=\"reflesh()\">\n      <i class=\"anticon anticon-reload\" style=\"color: #2db7f5\"></i>\n    </button>\n    <button style=\"margin: 4px;\" nz-button (click)=\"clear()\">\n        <i class=\"anticon anticon-delete\" style=\"color: #f50\"></i>\n    </button>\n  </div>\n</div>"

/***/ }),

/***/ "./src/app/components/dns-cache/dns-cache.component.ts":
/*!*************************************************************!*\
  !*** ./src/app/components/dns-cache/dns-cache.component.ts ***!
  \*************************************************************/
/*! exports provided: DnsCacheComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "DnsCacheComponent", function() { return DnsCacheComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _service_dns_cache_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../../service/dns-cache.service */ "./src/app/service/dns-cache.service.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var DnsCacheComponent = /** @class */ (function () {
    function DnsCacheComponent(service) {
        this.service = service;
    }
    DnsCacheComponent.prototype.ngOnInit = function () {
        this.reflesh();
    };
    DnsCacheComponent.prototype.reflesh = function () {
        var _this = this;
        this.service.getCache().subscribe(function (list) { return _this.cacheList = list; });
    };
    DnsCacheComponent.prototype.clear = function () {
        var _this = this;
        this.cacheList = [];
        this.service.clearCache().subscribe(function (_) {
            _this.reflesh();
        });
    };
    DnsCacheComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-dns-cache',
            template: __webpack_require__(/*! ./dns-cache.component.html */ "./src/app/components/dns-cache/dns-cache.component.html"),
            styles: [__webpack_require__(/*! ./dns-cache.component.css */ "./src/app/components/dns-cache/dns-cache.component.css")]
        }),
        __metadata("design:paramtypes", [_service_dns_cache_service__WEBPACK_IMPORTED_MODULE_1__["DnsCacheService"]])
    ], DnsCacheComponent);
    return DnsCacheComponent;
}());



/***/ }),

/***/ "./src/app/components/general/general.component.css":
/*!**********************************************************!*\
  !*** ./src/app/components/general/general.component.css ***!
  \**********************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ "./src/app/components/general/general.component.html":
/*!***********************************************************!*\
  !*** ./src/app/components/general/general.component.html ***!
  \***********************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<p>\n  general works!\n</p>\n"

/***/ }),

/***/ "./src/app/components/general/general.component.ts":
/*!*********************************************************!*\
  !*** ./src/app/components/general/general.component.ts ***!
  \*********************************************************/
/*! exports provided: GeneralComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "GeneralComponent", function() { return GeneralComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _service_general_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../../service/general.service */ "./src/app/service/general.service.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var GeneralComponent = /** @class */ (function () {
    function GeneralComponent(service) {
        this.service = service;
    }
    GeneralComponent.prototype.ngOnInit = function () {
    };
    GeneralComponent.prototype.shutdown = function () {
        this.service.shutdown().subscribe();
    };
    GeneralComponent.prototype.reload = function () {
        this.service.reload().subscribe();
    };
    GeneralComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-general',
            template: __webpack_require__(/*! ./general.component.html */ "./src/app/components/general/general.component.html"),
            styles: [__webpack_require__(/*! ./general.component.css */ "./src/app/components/general/general.component.css")]
        }),
        __metadata("design:paramtypes", [_service_general_service__WEBPACK_IMPORTED_MODULE_1__["GeneralService"]])
    ], GeneralComponent);
    return GeneralComponent;
}());



/***/ }),

/***/ "./src/app/components/menus/menus.component.css":
/*!******************************************************!*\
  !*** ./src/app/components/menus/menus.component.css ***!
  \******************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ".st-sider-menus{\n    height: 100vh;\n}\n.st-switch {\n    padding-top: 0.7em;\n    padding-bottom: 0.7em;\n}\n.st-switch div span {\n    margin-left: 1em;\n}"

/***/ }),

/***/ "./src/app/components/menus/menus.component.html":
/*!*******************************************************!*\
  !*** ./src/app/components/menus/menus.component.html ***!
  \*******************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<ul nz-menu class=\"st-sider-menus\" style=\"overflow: scroll;\">\n  <div (click)=\"githubHome()\" style=\"cursor: pointer;\">\n    <img src=\"/assets/menus_logo.png\" width=\"40px\" height=\"40px\" \n    style=\"float: left;margin: 5px;margin-left: 20px\"/>\n    <h3 style=\"height: 50px;line-height: 50px;margin: 5px;margin-left: 80px;margin-bottom: 0px\">Shuttle</h3>\n  </div>\n  <nz-divider style=\"margin-top: 10px;margin-bottom: 5px\"></nz-divider>\n  <li nz-menu-group>\n    <span title>Speed</span>\n    <div nz-row style=\"text-align: center;margin-bottom: 5px\">\n      <nz-tag [nzColor]=\"'blue'\">\n        <i class=\"anticon anticon-arrow-up\"></i>\n        {{speed.up_speed}}\n      </nz-tag>\n      <nz-tag [nzColor]=\"'blue'\">\n        {{speed.down_speed}}\n        <i class=\"anticon anticon-arrow-down\"></i>\n      </nz-tag>\n    </div>\n  </li>\n  <li nz-menu-group>\n    <span title>General</span>\n    <ul>\n      <li nz-menu-item routerLink=\"/servers\" \n      routerLinkActive #rla1_1=\"routerLinkActive\" \n      [nzSelected]=\"rla1_1.isActive\">Servers</li>\n      <li nz-menu-item routerLink=\"/mitm\" \n      routerLinkActive #rla1_2=\"routerLinkActive\" \n      [nzSelected]=\"rla1_2.isActive\">MITM</li>\n    </ul>\n  </li>\n  <li nz-menu-group>\n    <span title>HTTP Records</span>\n    <ul>\n      <li nz-menu-item routerLink=\"/records\" \n      routerLinkActive #rla2=\"routerLinkActive\" \n      [nzSelected]=\"rla2.isActive\">Records</li>\n    </ul>\n  </li>\n  <li nz-menu-group>\n    <span title>DNS</span>\n    <ul>\n      <li nz-menu-item routerLink=\"/dns-cache\"\n      routerLinkActive #rla3=\"routerLinkActive\" \n      [nzSelected]=\"rla3.isActive\">Cache</li>\n    </ul>\n  </li>\n  <li nz-menu-group>\n    <span title>Options</span>\n    <div style=\"padding-left: 10px; padding-right: 10px;\">\n      <nz-select style=\"width: 100%;\" [(ngModel)]=\"currentMode\"\n      (ngModelChange)=\"setMode($event)\">\n        <nz-option *ngFor=\"let mode of modeList\" \n        [nzLabel]=\"mode.label\" \n        [nzValue]=\"mode.value\"\n        ></nz-option>\n      </nz-select>\n    </div>\n    <div nz-row class=\"st-switch\">\n      <div nz-col nzOffset=\"4\">\n          <nz-switch [(ngModel)]=\"allow_dump\" (click)=\"dumpChange()\" ></nz-switch>\n          <span>Dump</span>\n      </div>\n    </div>\n    <div nz-row class=\"st-switch\">\n        <div nz-col nzOffset=\"4\">\n            <nz-switch [(ngModel)]=\"allow_mitm\" (click)=\"mitmChange()\" ></nz-switch>\n            <span>MITM</span>\n        </div>\n    </div>\n  </li>\n  <li nz-menu-group>\n    <div style=\"width: 100%; text-align: center;margin-bottom: 20px\">\n      <div style=\"width: 100%;text-align: center;padding: 5px;\">\n        <button nz-button nzType=\"default\" [nzSize]=\"size\" style=\"width: 130px;\" (click)=\"reload()\">\n            <i class=\"anticon anticon-reload\"></i>Reload\n        </button>\n      </div>\n      <div style=\"width: 100%;text-align: center;padding: 5px;\">\n        <button nz-button nzType=\"danger\" [nzSize]=\"size\" style=\"width: 130px\" (click)=\"shutdown()\">\n            <i class=\"anticon anticon-poweroff\"></i>Shutdown\n        </button>\n      </div>\n    </div>\n  </li>\n</ul>\n\n"

/***/ }),

/***/ "./src/app/components/menus/menus.component.ts":
/*!*****************************************************!*\
  !*** ./src/app/components/menus/menus.component.ts ***!
  \*****************************************************/
/*! exports provided: MenusComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "MenusComponent", function() { return MenusComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _service_dump_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../../service/dump.service */ "./src/app/service/dump.service.ts");
/* harmony import */ var _service_general_service__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../service/general.service */ "./src/app/service/general.service.ts");
/* harmony import */ var _modules_common_module__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../../modules/common.module */ "./src/app/modules/common.module.ts");
/* harmony import */ var _service_websocket_service__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../../service/websocket.service */ "./src/app/service/websocket.service.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var wsUpdateSpeed = _modules_common_module__WEBPACK_IMPORTED_MODULE_3__["WSHost"] + '/api/ws/speed';
var MenusComponent = /** @class */ (function () {
    function MenusComponent(dumpService, generalService, ws) {
        this.dumpService = dumpService;
        this.generalService = generalService;
        this.ws = ws;
        this.modeList = [
            { label: 'Rule Mode', value: 'RULE' },
            { label: 'Direct Mode', value: 'DIRECT' },
            { label: 'Remote Mode', value: 'REMOTE' },
            { label: 'Reject Mode', value: 'REJECT' }
        ];
    }
    MenusComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.dumpService.dumpStatus().subscribe(function (resp) {
            _this.allow_dump = resp.allow_dump;
            _this.allow_mitm = resp.allow_mitm;
        });
        this.generalService.getMode().subscribe(function (mode) { return _this.currentMode = mode; });
        this.speed = { up_speed: '0B/s', down_speed: '0B/s' };
        this.generalService.speed();
        // this.secondsCounter = interval(1000);
        // this.secondsCounter.subscribe(_ => this.getSpeed());
        this.ws.InitWebSocket(wsUpdateSpeed).subscribe(function (dataStr) {
            _this.speed = JSON.parse(dataStr);
        }, function (error) { return console.error(error); }, function () { return console.log('ws close!'); });
    };
    MenusComponent.prototype.dumpChange = function () {
        var _this = this;
        this.dumpService.allowDump(this.allow_dump).subscribe(function (resp) {
            _this.allow_dump = resp.allow_dump;
            _this.allow_mitm = resp.allow_mitm;
        });
    };
    MenusComponent.prototype.mitmChange = function () {
        var _this = this;
        this.dumpService.allowMitm(this.allow_mitm).subscribe(function (resp) {
            _this.allow_dump = resp.allow_dump;
            _this.allow_mitm = resp.allow_mitm;
        });
    };
    MenusComponent.prototype.shutdown = function () {
        this.generalService.shutdown().subscribe();
    };
    MenusComponent.prototype.reload = function () {
        this.generalService.reload().subscribe();
    };
    MenusComponent.prototype.githubHome = function () {
        window.open('https://github.com/sipt/shuttle');
    };
    MenusComponent.prototype.setMode = function (value) {
        var _this = this;
        this.generalService.setMode(value).subscribe(function (mode) { return _this.currentMode = mode; });
    };
    MenusComponent.prototype.getSpeed = function () {
        var _this = this;
        this.generalService.speed().subscribe(function (s) { return _this.speed = s; });
    };
    MenusComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-menus',
            template: __webpack_require__(/*! ./menus.component.html */ "./src/app/components/menus/menus.component.html"),
            styles: [__webpack_require__(/*! ./menus.component.css */ "./src/app/components/menus/menus.component.css")]
        }),
        __metadata("design:paramtypes", [_service_dump_service__WEBPACK_IMPORTED_MODULE_1__["DumpService"], _service_general_service__WEBPACK_IMPORTED_MODULE_2__["GeneralService"],
            _service_websocket_service__WEBPACK_IMPORTED_MODULE_4__["WebsocketService"]])
    ], MenusComponent);
    return MenusComponent;
}());



/***/ }),

/***/ "./src/app/components/mitm/mitm.component.css":
/*!****************************************************!*\
  !*** ./src/app/components/mitm/mitm.component.css ***!
  \****************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ "./src/app/components/mitm/mitm.component.html":
/*!*****************************************************!*\
  !*** ./src/app/components/mitm/mitm.component.html ***!
  \*****************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<nz-card style=\"width:300px;\" nzTitle=\"Certificate\">\n  <div style=\"width: 100%;text-align: center;padding: 5px;\">\n    <button nz-button nzType=\"default\" style=\"width: 130px;\" (click)=\"generate()\">\n      <i class=\"anticon anticon-file-text\"></i>\n      Generate\n    </button>\n  </div>\n  <div style=\"width: 100%;text-align: center;padding: 5px;\">\n    <button nz-button nzType=\"default\" style=\"width: 130px;\" (click)=\"download()\">\n      <i class=\"anticon anticon-download\"></i>\n      Download\n    </button>\n  </div>\n</nz-card>"

/***/ }),

/***/ "./src/app/components/mitm/mitm.component.ts":
/*!***************************************************!*\
  !*** ./src/app/components/mitm/mitm.component.ts ***!
  \***************************************************/
/*! exports provided: MitmComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "MitmComponent", function() { return MitmComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _service_general_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../../service/general.service */ "./src/app/service/general.service.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var MitmComponent = /** @class */ (function () {
    function MitmComponent(generalService) {
        this.generalService = generalService;
    }
    MitmComponent.prototype.ngOnInit = function () {
    };
    MitmComponent.prototype.generate = function () {
        this.generalService.generateCert().subscribe();
    };
    MitmComponent.prototype.download = function () {
        this.generalService.downloadCert();
    };
    MitmComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-mitm',
            template: __webpack_require__(/*! ./mitm.component.html */ "./src/app/components/mitm/mitm.component.html"),
            styles: [__webpack_require__(/*! ./mitm.component.css */ "./src/app/components/mitm/mitm.component.css")]
        }),
        __metadata("design:paramtypes", [_service_general_service__WEBPACK_IMPORTED_MODULE_1__["GeneralService"]])
    ], MitmComponent);
    return MitmComponent;
}());



/***/ }),

/***/ "./src/app/components/records/records.component.css":
/*!**********************************************************!*\
  !*** ./src/app/components/records/records.component.css ***!
  \**********************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ".st-ext {\n    width: 100%;\n    background: #ffffff;\n}\n.st-ext-title {\n    background: #fafafa;\n    width: 100%;\n    height: 40px;\n}\n.st-ext-content {\n    height: 460px;\n}\n.st-dump-content {\n    padding-left: 16px;\n    padding-right: 16px;\n    overflow: scroll;\n}\n"

/***/ }),

/***/ "./src/app/components/records/records.component.html":
/*!***********************************************************!*\
  !*** ./src/app/components/records/records.component.html ***!
  \***********************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<div [ngStyle]=\"tableStyle\">\n  <nz-table #list [nzData]=\"records\"\n  nzSize=\"small\"\n  [nzShowPagination]=\"false\"\n  [nzFrontPagination]=\"false\"\n  [nzScroll]=\"tbodyScroll\">\n    <thead id=\"st-thead\">\n      <tr>\n        <th nzWidth=\"100px\">ID</th>\n        <th nzWidth=\"150px\">Time</th>\n        <th nzWidth=\"100px\">Up/Down</th>\n        <th nzWidth=\"350px\">Policy</th>\n        <th nzWidth=\"100px\">Protocol</th>\n        <th style=\"padding: 0px\">\n          <nz-input-group nzSuffixIcon=\"anticon anticon-search\">\n            <input type=\"text\" nz-input placeholder=\"URL\" [(ngModel)]=\"keyword\" (ngModelChange)=\"filter($event)\">\n          </nz-input-group>\n        </th>\n      </tr>\n    </thead>\n    <tbody>\n      <tr *ngFor=\"let record of list.data\" (click)=\"getDump(record.ID, record.Dumped)\">\n        <td>{{record.ID}}</td>\n        <td>{{record.Created | date:'MM-dd HH:mm:ss'}}</td>\n        <td>{{record.Up | capacity}}/{{record.Down | capacity}}</td>\n        <td>{{record.Proxy.Name}}({{record.Rule.Type}},{{record.Rule.Value}})</td>\n        <td>{{record.Protocol}}</td>\n        <td>\n          <strong>\n            <i class=\"anticon anticon-swap\" *ngIf=\"record.Status=='Active'\" style=\"color: orange\"></i>\n          </strong>\n          <i class=\"anticon anticon-check-circle\" *ngIf=\"record.Status=='Completed' && !record.Dumped\" style=\"color: #87d068\"></i>\n          <i class=\"anticon anticon-close-circle\" *ngIf=\"record.Status=='Reject'\" style=\"color: #f50\"></i>\n          <strong *ngIf=\"record.Status=='Completed' && record.Dumped\">\n            <i class=\"anticon anticon-download\" style=\"color: #2db7f5\"></i>\n          </strong>\n          {{record.URL}}\n        </td>\n      </tr>\n    </tbody>\n  </nz-table>\n</div>\n<!-- <nz-affix nzOffsetBottom=\"0\" >\n</nz-affix> -->\n<div class=\"st-ext\">\n  <div class=\"st-ext-title\" (click)=\"closeExt()\">\n    <div>\n      <button style=\"margin: 4px;\" nz-button (click)=\"reflesh()\">\n        <i class=\"anticon anticon-reload\" style=\"color: #2db7f5\"></i>\n      </button>\n      <button style=\"margin: 4px;\" nz-button (click)=\"clear()\">\n          <i class=\"anticon anticon-delete\" style=\"color: #f50\"></i>\n      </button>\n      <div style=\"float: right;\">\n          <button *ngIf=\"!extClosed\" style=\"margin: 4px;border: 0px; background: #fafafa;\" nz-button>\n              <i class=\"anticon anticon-down\" style=\"color: #2db7f5\"></i>\n          </button>\n      </div>\n    </div>\n  </div>\n  <div class=\"st-ext-content\" *ngIf=\"!extClosed\">\n    <nz-tabset [nzSize]=\"small\" style=\"width: 100%;\">\n      <nz-tab nzTitle=\"Request-Header\">\n        <div class=\"st-dump-content\" style=\"height: 399px;\" [innerHTML]=\"dump.ReqHeader| nl2br\">\n        </div>\n      </nz-tab>\n      <nz-tab nzTitle=\"Request-Body\">\n        <div style=\"padding-left: 16px;width: 300px\">\n          <div *ngIf=\"dump.ReqBody\">\n            <nz-input-group nzSearch [nzSuffix]=\"suffixIconButton\">\n              <input type=\"text\" nz-input placeholder=\"input filename\" [(ngModel)]=\"reqFile\" >\n            </nz-input-group>\n            <ng-template #suffixIconButton>\n              <button nz-button nzType=\"primary\" nzSearch (click)=\"download('request')\">\n                <strong><i class=\"anticon anticon-download\"></i></strong>\n              </button>\n            </ng-template>\n          </div>\n        </div>\n        <div class=\"st-dump-content\" style=\"height: 367px;\" [innerHTML]=\"dump.ReqBody| nl2br\" *ngIf=\"dump.ReqBody\">\n        </div>\n        <div *ngIf=\"!dump.ReqBody\" style=\"text-align: center;height: 367px;\">\n          <div style=\"font-size: 5em;\">\n            <i class=\"anticon anticon-file\"></i>\n          </div>\n          Body is empty.\n        </div>\n      </nz-tab>\n      <nz-tab nzTitle=\"Response-Header\">\n        <div class=\"st-dump-content\" style=\"height: 399px;\" [innerHTML]=\"dump.RespHeader| nl2br\">\n        </div>\n      </nz-tab>\n      <nz-tab nzTitle=\"Response-Body\">\n        <div style=\"padding-left: 16px;width: 300px\">\n          <div *ngIf=\"dump.RespBody\">\n            <nz-input-group nzSearch [nzSuffix]=\"suffixIconButton\">\n              <input type=\"text\" nz-input placeholder=\"input filename\" [(ngModel)]=\"respFile\" >\n            </nz-input-group>\n            <ng-template #suffixIconButton>\n              <button nz-button nzType=\"primary\" nzSearch (click)=\"download('response')\">\n                <strong><i class=\"anticon anticon-download\"></i></strong>\n              </button>\n            </ng-template>\n          </div>\n        </div>\n        <div *ngIf=\"dump.RespBody!=='large body'\" class=\"st-dump-content\" style=\"height: 367px;\" [innerHTML]=\"dump.RespBody| html2text | nl2br\">\n        </div>\n        <div *ngIf=\"dump.RespBody==='large body'\" style=\"text-align: center;height: 367px;\">\n          <div style=\"font-size: 5em;\">\n            <i class=\"anticon anticon-file-unknown\"></i>\n          </div>\n          Large body, you can download it!\n        </div>\n      </nz-tab>\n    </nz-tabset>\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/components/records/records.component.ts":
/*!*********************************************************!*\
  !*** ./src/app/components/records/records.component.ts ***!
  \*********************************************************/
/*! exports provided: RecordsComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "RecordsComponent", function() { return RecordsComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _service_records_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../../service/records.service */ "./src/app/service/records.service.ts");
/* harmony import */ var _service_websocket_service__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../service/websocket.service */ "./src/app/service/websocket.service.ts");
/* harmony import */ var _modules_records_records_module__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../../modules/records/records.module */ "./src/app/modules/records/records.module.ts");
/* harmony import */ var _modules_common_module__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../../modules/common.module */ "./src/app/modules/common.module.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var wsUpdateRecords = _modules_common_module__WEBPACK_IMPORTED_MODULE_4__["WSHost"] + '/api/ws/records';
var largeFileURL = _modules_common_module__WEBPACK_IMPORTED_MODULE_4__["Host"] + '/api/dump/large';
var RecordsComponent = /** @class */ (function () {
    function RecordsComponent(service, ws) {
        this.service = service;
        this.ws = ws;
    }
    RecordsComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.reflesh();
        this.dump = new _modules_records_records_module__WEBPACK_IMPORTED_MODULE_3__["DumpModule"];
        this.extClosed = true;
        this.keyword = '';
        this.tbodyScroll = { y: 'calc(100vh - 80px)' };
        this.tableStyle = { height: 'calc(100vh - 40px)' };
        this.ws.InitWebSocket(wsUpdateRecords).subscribe(function (dataStr) {
            var data = JSON.parse(dataStr);
            if (data.Op === 4) {
                _this.allRecords = _this.allRecords.concat([{
                        ID: data.Value.ID,
                        Protocol: data.Value.Protocol,
                        Created: data.Value.Created,
                        Proxy: data.Value.Proxy,
                        Rule: data.Value.Rule,
                        Status: data.Value.Status,
                        Up: data.Value.Up,
                        Down: data.Value.Down,
                        URL: data.Value.URL,
                        Dumped: data.Value.Dumped,
                    }]);
                _this.filter();
            }
            else {
                _this.records.forEach(function (v, i) {
                    if (v.ID === data.ID) {
                        switch (data.Op) {
                            case 2:// 上传流量
                                _this.records[i].Up += data.Value;
                                break;
                            case 3:// 下载流量
                                _this.records[i].Down += data.Value;
                                break;
                            case 1:// 修改状态
                                _this.records[i].Status = data.Value;
                                break;
                            default:
                                break;
                        }
                    }
                });
            }
        }, function (error) { return console.error(error); }, function () { return console.log('ws close!'); });
        // this.secondsCounter = interval(1000);
        // this.secondsCounter.subscribe(_ => this.reflesh());
    };
    RecordsComponent.prototype.reflesh = function () {
        var _this = this;
        this.service.getCache().subscribe(function (list) {
            _this.allRecords = list;
            _this.filter();
        });
    };
    RecordsComponent.prototype.clear = function () {
        var _this = this;
        this.records = [];
        this.allRecords = [];
        this.service.clearCache().subscribe(function (_) { return _this.reflesh(); });
    };
    RecordsComponent.prototype.getDump = function (id, dumped) {
        var _this = this;
        if (dumped) {
            this.service.getDumpData(id).subscribe(function (d) {
                _this.dump = d;
            });
            this.openExt();
            this.id = id;
            this.reqFile = '';
            this.respFile = '';
        }
        else {
            this.dump = new _modules_records_records_module__WEBPACK_IMPORTED_MODULE_3__["DumpModule"];
        }
    };
    RecordsComponent.prototype.openExt = function () {
        if (this.extClosed) {
            this.extClosed = false;
            this.tbodyScroll = { y: 'calc(100vh - 540px)' };
            this.tableStyle = { height: 'calc(100vh - 500px)' };
        }
    };
    RecordsComponent.prototype.closeExt = function () {
        if (!this.extClosed) {
            this.extClosed = true;
            this.tbodyScroll = { y: 'calc(100vh - 80px)' };
            this.tableStyle = { height: 'calc(100vh - 40px)' };
        }
    };
    RecordsComponent.prototype.filter = function () {
        var _this = this;
        if (this.keyword !== '') {
            this.records = this.allRecords.filter(function (v) { return v.URL.indexOf(_this.keyword) >= 0; });
        }
        else {
            this.records = this.allRecords;
        }
    };
    RecordsComponent.prototype.download = function (dumpType) {
        var url = largeFileURL + '/' + this.id + '?file_name=';
        if (dumpType === 'request') {
            url += this.reqFile;
        }
        else {
            url += this.respFile;
        }
        url += '&dump_type=' + dumpType;
        window.open(url);
    };
    RecordsComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-records',
            template: __webpack_require__(/*! ./records.component.html */ "./src/app/components/records/records.component.html"),
            styles: [__webpack_require__(/*! ./records.component.css */ "./src/app/components/records/records.component.css")]
        }),
        __metadata("design:paramtypes", [_service_records_service__WEBPACK_IMPORTED_MODULE_1__["RecordsService"],
            _service_websocket_service__WEBPACK_IMPORTED_MODULE_2__["WebsocketService"]])
    ], RecordsComponent);
    return RecordsComponent;
}());



/***/ }),

/***/ "./src/app/components/server/server.component.css":
/*!********************************************************!*\
  !*** ./src/app/components/server/server.component.css ***!
  \********************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ "./src/app/components/server/server.component.html":
/*!*********************************************************!*\
  !*** ./src/app/components/server/server.component.html ***!
  \*********************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<div *ngFor=\"let group of list\" style=\"width: 300px;float:left;margin:0.7em\">\n  <nz-table\n    #servers\n    nzSize=\"small\"\n    [nzData]=\"group.servers\"\n    [nzShowPagination]=\"false\"\n    [nzFrontPagination]=\"false\">\n    <thead>\n      <tr>\n        <th>{{group.name}}({{group.select_type}})\n          <i class=\"anticon anticon-reload\"\n          *ngIf=\"group.select_type == 'rtt'\" \n          style=\"color: #2db7f5; float: right;\"\n          (click)=\"refreshSelect(group.name)\"></i>\n        </th>\n      </tr>\n    </thead>\n    <tbody>\n      <tr *ngFor=\"let server of servers.data\" (click)=\"select(group.name, server.name)\">\n        <td>{{server.name}}\n          <nz-tag *ngIf=\"server.rtt\">{{server.rtt}}</nz-tag>\n          <i *ngIf=\"server.selected\" class=\"anticon anticon-check-circle\" style=\"color: #87d068; float: right;\"></i>\n        </td>\n      </tr>\n    </tbody>\n  </nz-table>\n</div>\n"

/***/ }),

/***/ "./src/app/components/server/server.component.ts":
/*!*******************************************************!*\
  !*** ./src/app/components/server/server.component.ts ***!
  \*******************************************************/
/*! exports provided: ServerComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ServerComponent", function() { return ServerComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _service_server_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../../service/server.service */ "./src/app/service/server.service.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var ServerComponent = /** @class */ (function () {
    function ServerComponent(service) {
        this.service = service;
    }
    ServerComponent.prototype.ngOnInit = function () {
        this.refresh();
    };
    ServerComponent.prototype.refresh = function () {
        var _this = this;
        this.service.getServers().subscribe(function (list) {
            list.sort(function (x, y) { return y.servers.length - x.servers.length; });
            _this.list = list;
        });
    };
    ServerComponent.prototype.select = function (group, server) {
        var _this = this;
        for (var i = 0; i < this.list.length; i++) {
            var element = this.list[i];
            if (element.name === group && element.select_type === 'rtt') {
                return;
            }
        }
        this.service.selectServer(group, server).subscribe(function (_) { return _this.refresh(); });
    };
    ServerComponent.prototype.refreshSelect = function (group) {
        var _this = this;
        for (var i = 0; i < this.list.length; i++) {
            var element = this.list[i];
            if (element.name === group && element.select_type !== 'rtt') {
                return;
            }
        }
        this.service.refleshSelect(group).subscribe(function (_) { return _this.refresh(); });
    };
    ServerComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-server',
            template: __webpack_require__(/*! ./server.component.html */ "./src/app/components/server/server.component.html"),
            styles: [__webpack_require__(/*! ./server.component.css */ "./src/app/components/server/server.component.css")]
        }),
        __metadata("design:paramtypes", [_service_server_service__WEBPACK_IMPORTED_MODULE_1__["ServerService"]])
    ], ServerComponent);
    return ServerComponent;
}());



/***/ }),

/***/ "./src/app/modules/common.module.ts":
/*!******************************************!*\
  !*** ./src/app/modules/common.module.ts ***!
  \******************************************/
/*! exports provided: Response, Speed, Host, WSHost */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Response", function() { return Response; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Speed", function() { return Speed; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Host", function() { return Host; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "WSHost", function() { return WSHost; });
var Response = /** @class */ (function () {
    function Response() {
    }
    return Response;
}());

var Speed = /** @class */ (function () {
    function Speed() {
    }
    return Speed;
}());

var Host = 'http://localhost:8082';
var WSHost = 'ws://localhost:8082';
// export const WSHost = 'ws://' + document.location.host;
// export const Host = '';


/***/ }),

/***/ "./src/app/modules/records/records.module.ts":
/*!***************************************************!*\
  !*** ./src/app/modules/records/records.module.ts ***!
  \***************************************************/
/*! exports provided: RecordsModule, DumpModule */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "RecordsModule", function() { return RecordsModule; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "DumpModule", function() { return DumpModule; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common */ "./node_modules/@angular/common/fesm5/common.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};


var Rule = /** @class */ (function () {
    function Rule() {
    }
    Rule = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["NgModule"])({
            imports: [
                _angular_common__WEBPACK_IMPORTED_MODULE_1__["CommonModule"]
            ],
            declarations: []
        })
    ], Rule);
    return Rule;
}());
var Proxy = /** @class */ (function () {
    function Proxy() {
    }
    return Proxy;
}());
var RecordsModule = /** @class */ (function () {
    function RecordsModule() {
    }
    return RecordsModule;
}());

var DumpModule = /** @class */ (function () {
    function DumpModule() {
    }
    return DumpModule;
}());



/***/ }),

/***/ "./src/app/pipes/capacity.pipe.ts":
/*!****************************************!*\
  !*** ./src/app/pipes/capacity.pipe.ts ***!
  \****************************************/
/*! exports provided: CapacityPipe */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "CapacityPipe", function() { return CapacityPipe; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};

var CapacityPipe = /** @class */ (function () {
    function CapacityPipe() {
    }
    CapacityPipe.prototype.transform = function (value, args) {
        var unit = 'B';
        var t = value;
        var n = Math.floor(t / 1024);
        if (n >= 1) {
            t = n;
            unit = 'KB';
            n = Math.floor(t / 1024);
            if (n >= 1) {
                t = n;
                unit = 'MB';
                n = Math.floor(t / 1024);
                if (n >= 1) {
                    t = n;
                    unit = 'GB';
                }
            }
        }
        return t + unit;
    };
    CapacityPipe = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Pipe"])({
            name: 'capacity'
        })
    ], CapacityPipe);
    return CapacityPipe;
}());



/***/ }),

/***/ "./src/app/pipes/html2text.pipe.ts":
/*!*****************************************!*\
  !*** ./src/app/pipes/html2text.pipe.ts ***!
  \*****************************************/
/*! exports provided: Html2textPipe */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Html2textPipe", function() { return Html2textPipe; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};

var Html2textPipe = /** @class */ (function () {
    function Html2textPipe() {
    }
    Html2textPipe.prototype.transform = function (value, args) {
        if (typeof value !== 'string') {
            return value;
        }
        var textParsed = value.replace(/(?:&)/g, '&amp;')
            .replace(/(?:")/g, '&quot;')
            .replace(/(?:<)/g, '&lt;')
            .replace(/(?:>)/g, '&gt;')
            .replace(/(?: )/g, '&nbsp;');
        return textParsed;
    };
    Html2textPipe = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Pipe"])({
            name: 'html2text'
        })
    ], Html2textPipe);
    return Html2textPipe;
}());



/***/ }),

/***/ "./src/app/pipes/ips-format.pipe.ts":
/*!******************************************!*\
  !*** ./src/app/pipes/ips-format.pipe.ts ***!
  \******************************************/
/*! exports provided: IpsFormatPipe */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "IpsFormatPipe", function() { return IpsFormatPipe; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};

var IpsFormatPipe = /** @class */ (function () {
    function IpsFormatPipe() {
    }
    IpsFormatPipe.prototype.transform = function (value, args) {
        var arr = value;
        var values = '';
        arr.forEach(function (v, i) {
            if (i > 0) {
                values += ',';
            }
            values += "<div>" + v + "</div>";
        });
        return values;
    };
    IpsFormatPipe = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Pipe"])({
            name: 'ipsFormat'
        })
    ], IpsFormatPipe);
    return IpsFormatPipe;
}());



/***/ }),

/***/ "./src/app/pipes/nl2br.pipe.ts":
/*!*************************************!*\
  !*** ./src/app/pipes/nl2br.pipe.ts ***!
  \*************************************/
/*! exports provided: Nl2BrPipe */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Nl2BrPipe", function() { return Nl2BrPipe; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_platform_browser__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/platform-browser */ "./node_modules/@angular/platform-browser/fesm5/platform-browser.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var Nl2BrPipe = /** @class */ (function () {
    function Nl2BrPipe(sanitizer) {
        this.sanitizer = sanitizer;
    }
    Nl2BrPipe.prototype.transform = function (value, sanitizeBeforehand) {
        if (typeof value !== 'string') {
            return value;
        }
        var result;
        var textParsed = value.replace(/(?:\r\n|\r|\n)/g, '<br />');
        if (!_angular_core__WEBPACK_IMPORTED_MODULE_0__["VERSION"] || _angular_core__WEBPACK_IMPORTED_MODULE_0__["VERSION"].major === '2') {
            result = this.sanitizer.bypassSecurityTrustHtml(textParsed);
        }
        else if (sanitizeBeforehand) {
            result = this.sanitizer.sanitize(_angular_core__WEBPACK_IMPORTED_MODULE_0__["SecurityContext"].HTML, textParsed);
        }
        else {
            result = textParsed;
        }
        return result;
    };
    Nl2BrPipe = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Pipe"])({
            name: 'nl2br'
        }),
        __metadata("design:paramtypes", [_angular_platform_browser__WEBPACK_IMPORTED_MODULE_1__["DomSanitizer"]])
    ], Nl2BrPipe);
    return Nl2BrPipe;
}());



/***/ }),

/***/ "./src/app/service/dns-cache.service.ts":
/*!**********************************************!*\
  !*** ./src/app/service/dns-cache.service.ts ***!
  \**********************************************/
/*! exports provided: DnsCacheService */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "DnsCacheService", function() { return DnsCacheService; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var _modules_common_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../modules/common.module */ "./src/app/modules/common.module.ts");
/* harmony import */ var rxjs__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! rxjs */ "./node_modules/rxjs/_esm5/index.js");
/* harmony import */ var rxjs_operators__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! rxjs/operators */ "./node_modules/rxjs/_esm5/operators/index.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var DnsCacheUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/dns';
var DnsCacheService = /** @class */ (function () {
    function DnsCacheService(http) {
        this.http = http;
    }
    DnsCacheService.prototype.getCache = function () {
        return this.http.get(DnsCacheUrl)
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getCache', {
            code: 1,
            message: '',
            data: [],
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    DnsCacheService.prototype.clearCache = function () {
        return this.http.delete(DnsCacheUrl).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getCache', {
            code: 1,
            message: '',
            data: [],
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) {
            var r = resp;
            if (r.code === 1) {
                console.error(r.message);
            }
            return 1;
        }));
    };
    /**
     * Handle Http operation that failed.
     * Let the app continue.
     * @param operation - name of the operation that failed
     * @param result - optional value to return as the observable result
     */
    DnsCacheService.prototype.handleError = function (operation, result) {
        if (operation === void 0) { operation = 'operation'; }
        return function (error) {
            // TODO: send the error to remote logging infrastructure
            console.error(error); // log to console instead
            // TODO: better job of transforming error for user consumption
            // this.log(`${operation} failed: ${error.message}`);
            // Let the app keep running by returning an empty result.
            return Object(rxjs__WEBPACK_IMPORTED_MODULE_3__["of"])(result);
        };
    };
    DnsCacheService = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Injectable"])({
            providedIn: 'root'
        }),
        __metadata("design:paramtypes", [_angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpClient"]])
    ], DnsCacheService);
    return DnsCacheService;
}());



/***/ }),

/***/ "./src/app/service/dump.service.ts":
/*!*****************************************!*\
  !*** ./src/app/service/dump.service.ts ***!
  \*****************************************/
/*! exports provided: DumpService */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "DumpService", function() { return DumpService; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var _modules_common_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../modules/common.module */ "./src/app/modules/common.module.ts");
/* harmony import */ var rxjs__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! rxjs */ "./node_modules/rxjs/_esm5/index.js");
/* harmony import */ var rxjs_operators__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! rxjs/operators */ "./node_modules/rxjs/_esm5/operators/index.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var AllowDumpUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/dump/allow';
var DumpService = /** @class */ (function () {
    function DumpService(http) {
        this.http = http;
    }
    DumpService.prototype.dumpStatus = function () {
        return this.http.get(AllowDumpUrl)
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getCache', {
            code: 1,
            message: '',
            data: { allow_dump: false, allow_mitm: false }
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    DumpService.prototype.allowDump = function (allow) {
        var headers = new _angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpHeaders"]().set('Content-Type', 'application/x-www-form-urlencoded; charset=utf-8');
        return this.http.post(AllowDumpUrl, 'allow_dump=' + allow, { headers: headers })
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getCache', {
            code: 1,
            message: '',
            data: { allow_dump: false, allow_mitm: false }
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    DumpService.prototype.allowMitm = function (allow) {
        var headers = new _angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpHeaders"]().set('Content-Type', 'application/x-www-form-urlencoded; charset=utf-8');
        return this.http.post(AllowDumpUrl, 'allow_mitm=' + allow, { headers: headers })
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getCache', {
            code: 1,
            message: '',
            data: { allow_dump: false, allow_mitm: false }
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    /**
     * Handle Http operation that failed.
     * Let the app continue.
     * @param operation - name of the operation that failed
     * @param result - optional value to return as the observable result
     */
    DumpService.prototype.handleError = function (operation, result) {
        if (operation === void 0) { operation = 'operation'; }
        return function (error) {
            // TODO: send the error to remote logging infrastructure
            console.error(error); // log to console instead
            // TODO: better job of transforming error for user consumption
            // this.log(`${operation} failed: ${error.message}`);
            // Let the app keep running by returning an empty result.
            return Object(rxjs__WEBPACK_IMPORTED_MODULE_3__["of"])(result);
        };
    };
    DumpService = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Injectable"])({
            providedIn: 'root'
        }),
        __metadata("design:paramtypes", [_angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpClient"]])
    ], DumpService);
    return DumpService;
}());

var DumpStatus = /** @class */ (function () {
    function DumpStatus() {
    }
    return DumpStatus;
}());


/***/ }),

/***/ "./src/app/service/general.service.ts":
/*!********************************************!*\
  !*** ./src/app/service/general.service.ts ***!
  \********************************************/
/*! exports provided: GeneralService */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "GeneralService", function() { return GeneralService; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var _modules_common_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../modules/common.module */ "./src/app/modules/common.module.ts");
/* harmony import */ var rxjs__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! rxjs */ "./node_modules/rxjs/_esm5/index.js");
/* harmony import */ var rxjs_operators__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! rxjs/operators */ "./node_modules/rxjs/_esm5/operators/index.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var shutdownUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/shutdown';
var reloadUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/reload';
var certUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/cert';
var modeUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/mode';
var speedUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/speed';
var GeneralService = /** @class */ (function () {
    function GeneralService(http) {
        this.http = http;
    }
    GeneralService.prototype.shutdown = function () {
        return this.http.post(shutdownUrl, {})
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('shutdown', {
            code: 1,
            message: '',
            data: { allow_dump: false, allow_mitm: false }
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    GeneralService.prototype.reload = function () {
        return this.http.post(reloadUrl, {})
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('reload', {
            code: 1,
            message: '',
            data: {}
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    GeneralService.prototype.generateCert = function () {
        return this.http.post(certUrl, {})
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('generateCert', {
            code: 1,
            message: '',
            data: {}
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    GeneralService.prototype.downloadCert = function () {
        window.open(certUrl);
    };
    GeneralService.prototype.getMode = function () {
        return this.http.get(modeUrl)
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getMode', {
            code: 1,
            message: '',
            data: {}
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    GeneralService.prototype.setMode = function (mode) {
        return this.http.post(modeUrl + '/' + mode, {})
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getMode', {
            code: 1,
            message: '',
            data: {}
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    GeneralService.prototype.speed = function () {
        return this.http.get(speedUrl)
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('speed', {
            code: 1,
            message: '',
            data: {}
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    /**
     * Handle Http operation that failed.
     * Let the app continue.
     * @param operation - name of the operation that failed
     * @param result - optional value to return as the observable result
     */
    GeneralService.prototype.handleError = function (operation, result) {
        if (operation === void 0) { operation = 'operation'; }
        return function (error) {
            // TODO: send the error to remote logging infrastructure
            console.error(error); // log to console instead
            // TODO: better job of transforming error for user consumption
            // this.log(`${operation} failed: ${error.message}`);
            // Let the app keep running by returning an empty result.
            return Object(rxjs__WEBPACK_IMPORTED_MODULE_3__["of"])(result);
        };
    };
    GeneralService = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Injectable"])({
            providedIn: 'root'
        }),
        __metadata("design:paramtypes", [_angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpClient"]])
    ], GeneralService);
    return GeneralService;
}());



/***/ }),

/***/ "./src/app/service/records.service.ts":
/*!********************************************!*\
  !*** ./src/app/service/records.service.ts ***!
  \********************************************/
/*! exports provided: RecordsService */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "RecordsService", function() { return RecordsService; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var _modules_common_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../modules/common.module */ "./src/app/modules/common.module.ts");
/* harmony import */ var _modules_records_records_module__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../modules/records/records.module */ "./src/app/modules/records/records.module.ts");
/* harmony import */ var rxjs__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! rxjs */ "./node_modules/rxjs/_esm5/index.js");
/* harmony import */ var rxjs_operators__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! rxjs/operators */ "./node_modules/rxjs/_esm5/operators/index.js");
/* harmony import */ var _utils_utils__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! ../utils/utils */ "./src/app/utils/utils.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};







var RecordsUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/records';
var DumpDataUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/dump/data/';
var utf8 = new _utils_utils__WEBPACK_IMPORTED_MODULE_6__["UTF8"]();
var RecordsService = /** @class */ (function () {
    function RecordsService(http) {
        this.http = http;
    }
    RecordsService.prototype.getCache = function () {
        return this.http.get(RecordsUrl)
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_5__["catchError"])(this.handleError('getCache', {
            code: 1,
            message: '',
            data: [],
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_5__["map"])(function (resp) { return resp.data; }));
    };
    RecordsService.prototype.clearCache = function () {
        return this.http.delete(RecordsUrl).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_5__["catchError"])(this.handleError('getCache', {
            code: 1,
            message: '',
            data: [],
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_5__["map"])(function (resp) {
            var r = resp;
            if (r.code === 1) {
                console.error(r.message);
            }
            return 1;
        }));
    };
    RecordsService.prototype.getDumpData = function (id) {
        var headers = new _angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpHeaders"]().set('Content-Type', 'application/json; charset=utf-8');
        return this.http.get(DumpDataUrl + id, { headers: headers }).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_5__["catchError"])(this.handleError('getDumpData', {
            code: 1,
            message: '',
            data: {},
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_5__["map"])(function (resp) {
            var r = resp;
            if (r.code === 1) {
                console.error(r.message);
            }
            var dump = new _modules_records_records_module__WEBPACK_IMPORTED_MODULE_3__["DumpModule"]();
            if (r.data.ReqHeader !== '') {
                dump.ReqHeader = utf8.decode(atob(r.data.ReqHeader));
            }
            if (r.data.ReqBody !== '') {
                dump.ReqBody = utf8.decode(atob(r.data.ReqBody));
            }
            if (r.data.RespHeader !== '') {
                dump.RespHeader = utf8.decode(atob(r.data.RespHeader));
            }
            if (r.data.RespBody !== '') {
                dump.RespBody = utf8.decode(atob(r.data.RespBody));
            }
            return dump;
        }));
    };
    /**
     * Handle Http operation that failed.
     * Let the app continue.
     * @param operation - name of the operation that failed
     * @param result - optional value to return as the observable result
     */
    RecordsService.prototype.handleError = function (operation, result) {
        if (operation === void 0) { operation = 'operation'; }
        return function (error) {
            // TODO: send the error to remote logging infrastructure
            console.error(error); // log to console instead
            // TODO: better job of transforming error for user consumption
            // this.log(`${operation} failed: ${error.message}`);
            // Let the app keep running by returning an empty result.
            return Object(rxjs__WEBPACK_IMPORTED_MODULE_4__["of"])(result);
        };
    };
    RecordsService = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Injectable"])({
            providedIn: 'root'
        }),
        __metadata("design:paramtypes", [_angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpClient"]])
    ], RecordsService);
    return RecordsService;
}());



/***/ }),

/***/ "./src/app/service/server.service.ts":
/*!*******************************************!*\
  !*** ./src/app/service/server.service.ts ***!
  \*******************************************/
/*! exports provided: ServerService */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ServerService", function() { return ServerService; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var _modules_common_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../modules/common.module */ "./src/app/modules/common.module.ts");
/* harmony import */ var rxjs__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! rxjs */ "./node_modules/rxjs/_esm5/index.js");
/* harmony import */ var rxjs_operators__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! rxjs/operators */ "./node_modules/rxjs/_esm5/operators/index.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var ServersUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/servers';
var SelectServerUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/server/select';
var RefreshSelectUrl = _modules_common_module__WEBPACK_IMPORTED_MODULE_2__["Host"] + '/api/server/select/refresh';
var ServerService = /** @class */ (function () {
    function ServerService(http) {
        this.http = http;
    }
    ServerService.prototype.getServers = function () {
        return this.http.get(ServersUrl)
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('getServers', {
            code: 1,
            message: '',
            data: [],
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    ServerService.prototype.selectServer = function (group, server) {
        var headers = new _angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpHeaders"]().set('Content-Type', 'application/x-www-form-urlencoded; charset=utf-8');
        return this.http.post(SelectServerUrl, 'group=' + group + '&server=' + server, { headers: headers })
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('selectServer', {
            code: 1,
            message: '',
            data: {}
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    ServerService.prototype.refleshSelect = function (group) {
        var headers = new _angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpHeaders"]().set('Content-Type', 'application/x-www-form-urlencoded; charset=utf-8');
        return this.http.post(RefreshSelectUrl, 'group=' + group, { headers: headers })
            .pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["catchError"])(this.handleError('refleshSelect', {
            code: 1,
            message: '',
            data: {}
        }))).pipe(Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_4__["map"])(function (resp) { return resp.data; }));
    };
    /**
     * Handle Http operation that failed.
     * Let the app continue.
     * @param operation - name of the operation that failed
     * @param result - optional value to return as the observable result
     */
    ServerService.prototype.handleError = function (operation, result) {
        if (operation === void 0) { operation = 'operation'; }
        return function (error) {
            // TODO: send the error to remote logging infrastructure
            console.error(error); // log to console instead
            // TODO: better job of transforming error for user consumption
            // this.log(`${operation} failed: ${error.message}`);
            // Let the app keep running by returning an empty result.
            return Object(rxjs__WEBPACK_IMPORTED_MODULE_3__["of"])(result);
        };
    };
    ServerService = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Injectable"])({
            providedIn: 'root'
        }),
        __metadata("design:paramtypes", [_angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpClient"]])
    ], ServerService);
    return ServerService;
}());



/***/ }),

/***/ "./src/app/service/websocket.service.ts":
/*!**********************************************!*\
  !*** ./src/app/service/websocket.service.ts ***!
  \**********************************************/
/*! exports provided: WebsocketService */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "WebsocketService", function() { return WebsocketService; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _node_modules_rxjs__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../../../node_modules/rxjs */ "./node_modules/rxjs/_esm5/index.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var WebsocketService = /** @class */ (function () {
    function WebsocketService() {
    }
    WebsocketService.prototype.InitWebSocket = function (url) {
        var _this = this;
        this.ws = new WebSocket(url);
        return new _node_modules_rxjs__WEBPACK_IMPORTED_MODULE_1__["Observable"](function (observer) {
            _this.ws.onmessage = function (event) { return observer.next(event.data); };
            _this.ws.onerror = function (event) { return observer.next(event); };
            _this.ws.onclose = function (event) { return observer.complete(); };
        });
    };
    WebsocketService.prototype.sendMessage = function (data) {
        this.ws.send(data);
    };
    WebsocketService = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Injectable"])({
            providedIn: 'root'
        }),
        __metadata("design:paramtypes", [])
    ], WebsocketService);
    return WebsocketService;
}());



/***/ }),

/***/ "./src/app/utils/utils.ts":
/*!********************************!*\
  !*** ./src/app/utils/utils.ts ***!
  \********************************/
/*! exports provided: UTF8 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "UTF8", function() { return UTF8; });
var UTF8 = /** @class */ (function () {
    function UTF8() {
    }
    /**
           * Encode multi-byte Unicode string into utf-8 multiple single-byte characters
           * (BMP / basic multilingual plane only)
           *
           * Chars in range U+0080 - U+07FF are encoded in 2 chars, U+0800 - U+FFFF in 3 chars
           *
           * @param {String} strUni Unicode string to be encoded as UTF-8
           * @returns {String} encoded string
           */
    UTF8.prototype.encode = function (strUni) {
        // use regular expressions & String.replace callback function for better efficiency
        // than procedural approaches
        var strUtf = strUni.replace(/[\u0080-\u07ff]/g, // U+0080 - U+07FF => 2 bytes 110yyyyy, 10zzzzzz
        function (// U+0080 - U+07FF => 2 bytes 110yyyyy, 10zzzzzz
        c) {
            var cc = c.charCodeAt(0);
            return String.fromCharCode(0xc0 | cc >> 6, 0x80 | cc & 0x3f);
        })
            .replace(/[\u0800-\uffff]/g, // U+0800 - U+FFFF => 3 bytes 1110xxxx, 10yyyyyy, 10zzzzzz
        function (// U+0800 - U+FFFF => 3 bytes 1110xxxx, 10yyyyyy, 10zzzzzz
        c) {
            var cc = c.charCodeAt(0);
            return String.fromCharCode(0xe0 | cc >> 12, 0x80 | cc >> 6 & 0x3F, 0x80 | cc & 0x3f);
        });
        return strUtf;
    };
    /**
     * Decode utf-8 encoded string back into multi-byte Unicode characters
     *
     * @param {String} strUtf UTF-8 string to be decoded back to Unicode
     * @returns {String} decoded string
     */
    UTF8.prototype.decode = function (strUtf) {
        // note: decode 3-byte chars first as decoded 2-byte strings could appear to be 3-byte char!
        var strUni = strUtf.replace(/[\u00e0-\u00ef][\u0080-\u00bf][\u0080-\u00bf]/g, // 3-byte chars
        function (// 3-byte chars
        c) {
            var cc = ((c.charCodeAt(0) & 0x0f) << 12) | ((c.charCodeAt(1) & 0x3f) << 6) | (c.charCodeAt(2) & 0x3f);
            return String.fromCharCode(cc);
        })
            .replace(/[\u00c0-\u00df][\u0080-\u00bf]/g, // 2-byte chars
        function (// 2-byte chars
        c) {
            var cc = (c.charCodeAt(0) & 0x1f) << 6 | c.charCodeAt(1) & 0x3f;
            return String.fromCharCode(cc);
        });
        return strUni;
    };
    return UTF8;
}());



/***/ }),

/***/ "./src/environments/environment.ts":
/*!*****************************************!*\
  !*** ./src/environments/environment.ts ***!
  \*****************************************/
/*! exports provided: environment */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "environment", function() { return environment; });
// This file can be replaced during build by using the `fileReplacements` array.
// `ng build ---prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.
var environment = {
    production: false
};
/*
 * In development mode, for easier debugging, you can ignore zone related error
 * stack frames such as `zone.run`/`zoneDelegate.invokeTask` by importing the
 * below file. Don't forget to comment it out in production mode
 * because it will have a performance impact when errors are thrown
 */
// import 'zone.js/dist/zone-error';  // Included with Angular CLI.


/***/ }),

/***/ "./src/main.ts":
/*!*********************!*\
  !*** ./src/main.ts ***!
  \*********************/
/*! no exports provided */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_platform_browser_dynamic__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/platform-browser-dynamic */ "./node_modules/@angular/platform-browser-dynamic/fesm5/platform-browser-dynamic.js");
/* harmony import */ var _app_app_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./app/app.module */ "./src/app/app.module.ts");
/* harmony import */ var _environments_environment__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./environments/environment */ "./src/environments/environment.ts");




if (_environments_environment__WEBPACK_IMPORTED_MODULE_3__["environment"].production) {
    Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["enableProdMode"])();
}
Object(_angular_platform_browser_dynamic__WEBPACK_IMPORTED_MODULE_1__["platformBrowserDynamic"])().bootstrapModule(_app_app_module__WEBPACK_IMPORTED_MODULE_2__["AppModule"])
    .catch(function (err) { return console.log(err); });


/***/ }),

/***/ 0:
/*!***************************!*\
  !*** multi ./src/main.ts ***!
  \***************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__(/*! /Users/sipt/workspace/js/shuttle/src/main.ts */"./src/main.ts");


/***/ })

},[[0,"runtime","vendor"]]]);
//# sourceMappingURL=main.js.map