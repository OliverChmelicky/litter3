//Source which seems to be better https://www.techiediaries.com/angular-firebase/angular-9-firebase-authentication-email-google-and-password/

import {Injectable} from '@angular/core';
import * as firebase from 'firebase/app';
import {auth} from 'firebase/app';
import {AngularFireAuth} from "@angular/fire/auth";
import {UserService} from "../user/user.service";
import { Router } from '@angular/router';
import {MemberModel, SocietyModel, UserModel} from "../../models/user.model";
import {BehaviorSubject, throwError} from "rxjs";


@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private loggedIn = new BehaviorSubject<boolean>(false);

  constructor(
    private  afAuth: AngularFireAuth,
    private userService: UserService,
    private router: Router,
  ) {
    this.afAuth.authState.subscribe(user => {
      if (user) {
        localStorage.setItem('firebaseUser', JSON.stringify(user));

        user.getIdToken().then(token => localStorage.setItem('token', token))
          .catch(() => localStorage.setItem('token', null));
      } else {
        localStorage.setItem('firebaseUser', null);
        localStorage.setItem('token', null);
      }
    })
  }


  async login(email: string, password: string) {
    await this.afAuth.signInWithEmailAndPassword(email, password)
      .then(res => {
        localStorage.setItem('firebaseUser', JSON.stringify(res.user));
        res.user.getIdToken().then(token => {
          localStorage.setItem('token', token)
          this.loggedIn.next(true);
          this.router.navigate(['/me']);
        })
          .catch(() =>
            localStorage.setItem('token', null)
          );
      }).catch(() => {
        localStorage.setItem('firebaseUser', null);
        localStorage.setItem('token', null);
        return null;
      })
  }

  renewToken() {
    this.afAuth.currentUser.then(user => {
      user.getIdToken(true)
        .then((token) =>
          localStorage.setItem('token', token)
        ).catch(
        err => {
          console.log('error custom renew token inside ' + err);
          throwError(err)
        })
    })
      .catch(err => {
        console.log('error custom renew token first ' + err);
        throwError(err)
      })
  }

  async register(value) {
    await this.afAuth.createUserWithEmailAndPassword(value.email, value.password)
      .then(res => {
        const firebaseUser = res.user;
        localStorage.setItem('firebaseUser', JSON.stringify(res.user));
        this.userService.createUser({
          Id: '',
          FirstName: 'Olo',//value.FirstName,
          LastName: 'Chmelo',//value.LastName,
          Email: value.email,
          Uid: firebaseUser.uid,
          Avatar: '',
          CreatedAt: new Date()
        }).subscribe(
          () => {
            this.renewToken();
            this.loggedIn.next(true);
            this.router.navigate(['/me']);
          },
          err => throwError(err)
        );
      })
      .catch(err => {
          localStorage.setItem('firebaseUser', null);
          localStorage.setItem('token', null);
          throw err;
        }
      )
  }

  async sendPasswordResetEmail(passwordResetEmail: string) {
    return await this.afAuth.sendPasswordResetEmail(passwordResetEmail);
  }

  async logout() {
    await this.afAuth.signOut();
    this.loggedIn.next(false);
    this.router.navigate(['/login']);
    localStorage.setItem('firebaseUser', null);
    localStorage.setItem('token', null);
  }

  get isLoggedIn() {
    const isLogged = localStorage.getItem('firebaseUser')
    if (isLogged != null) {
      this.loggedIn.next(true)
      this.router.navigate(['/me']);
    } else {
      this.loggedIn.next(false)
      this.router.navigate(['/login']);
    }
    return this.loggedIn.asObservable();
  }

  async loginWithGoogle() {
    await this.afAuth.signInWithPopup(new auth.GoogleAuthProvider())
  }

  getToken(): string {
    return localStorage.getItem('token')
  }
}

//zdroj Č. 2
//
//https://angular-templates.io/tutorials/about/firebase-authentication-with-angular
// doGoogleLogin(){
//   return new Promise<any>((resolve, reject) => {
//     let provider = new firebase.auth.GoogleAuthProvider();
//     provider.addScope('profile');
//     provider.addScope('email');
//     this.afAuth.auth
//       .signInWithPopup(provider)
//       .then(res => {
//         resolve(res);
//       })
//   })
// }
//
// doFacebookLogin(){
//   return new Promise<any>((resolve, reject) => {
//     let provider = new firebase.auth.FacebookAuthProvider();
//     this.afAuth.auth
//       .signInWithPopup(provider)
//       .then(res => {
//         resolve(res);
//       }, err => {
//         console.log(err);
//         reject(err);
//       })
//   })
// }
//
// doRegister(value){
//   return new Promise<any>((resolve, reject) => {
//     firebase.auth().createUserWithEmailAndPassword(value.email, value.password)
//       .then(res => {
//         resolve(res);
//       }, err => reject(err))
//   })
// }
//END OF SOURCE
