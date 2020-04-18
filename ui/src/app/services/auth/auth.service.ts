//Source which seems to be better https://www.techiediaries.com/angular-firebase/angular-9-firebase-authentication-email-google-and-password/

import {Injectable} from '@angular/core';
import * as firebase from 'firebase/app';
import {auth} from 'firebase/app';
import {AngularFireAuth} from "@angular/fire/auth";
import {User} from 'firebase';
import {tokenNotExpired} from "angular2-jwt";


@Injectable({
  providedIn: 'root'
})
export class AuthService {
  user: User;
  token: string;

  constructor(
    public  afAuth: AngularFireAuth,
  ) {
    this.afAuth.authState.subscribe(user => {
      if (user) {
        this.user = user;
        localStorage.setItem('user', JSON.stringify(this.user));
      } else {
        localStorage.setItem('user', null);
      }
    })
  }


  async login(email: string, password: string) {
    await this.afAuth.auth.signInWithEmailAndPassword(email, password)
      .then(res => {
        this.user = res.user;
        localStorage.setItem('user', JSON.stringify(this.user));
        return res.user;
      }).catch(() => {
          this.user = null;
          localStorage.setItem('user', null);
          return null;
        }
      )
  }

  async register(email: string, password: string): Promise<firebase.auth.UserCredential> {
    return this.afAuth.auth.createUserWithEmailAndPassword(email, password)
      .then(res => {
        this.user = res.user;
        localStorage.setItem('user', JSON.stringify(this.user));
        return res.user;
      }).catch(() => {
          this.user = null;
          localStorage.setItem('user', null);
          return null;
        }
      )
  }

  async sendPasswordResetEmail(passwordResetEmail: string) {
    return await this.afAuth.auth.sendPasswordResetEmail(passwordResetEmail);
  }

  async logout() {
    await this.afAuth.auth.signOut();
    localStorage.removeItem('user');
  }

  get isLoggedIn(): boolean {
    const user = JSON.parse(localStorage.getItem('user'));
    return user !== null;
  }

  async loginWithGoogle() {
    await this.afAuth.auth.signInWithPopup(new auth.GoogleAuthProvider())
  }

  async getToken() {
    const user = JSON.parse(localStorage.getItem('user'));
    return await (user as firebase.User).getIdToken().then(val => this.token = val)
  }
  public isAuthenticated(): boolean {
    // get the token
    this.getToken();
    // return a boolean reflecting
    // whether or not the token is expired
    return tokenNotExpired(null, this.token);
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
