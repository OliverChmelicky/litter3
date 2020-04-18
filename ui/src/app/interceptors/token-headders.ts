// import { Injectable } from "@angular/core";
// import {
//   HttpRequest,
//   HttpHandler,
//   HttpEvent,
//   HttpInterceptor
// } from "@angular/common/http";
// import { Observable } from "rxjs";
// import {AuthService} from "../services/auth/auth.service";
//
// //Add check na to ze token vyprsal
// //https://medium.com/@ryanchenkie_40935/angular-authentication-using-the-http-client-and-http-interceptors-2f9d1540eb8
//
// @Injectable()
// export class TokenInterceptor implements HttpInterceptor {
//   constructor(public auth: AuthService) {}
//   intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
//
//     request = request.clone({
//       setHeaders: {
//         Authorization: `${this.auth.getToken()}`
//       }
//     });
//     return next.handle(request);
//   }
// }
