//Source which seems to be better https://www.techiediaries.com/angular-firebase/angular-9-firebase-authentication-email-google-and-password/

import {Injectable} from '@angular/core';
import * as firebase from 'firebase/app';
import {auth} from 'firebase/app';
import {AngularFireAuth} from "@angular/fire/auth";
import {UserService} from "../user/user.service";
import {MemberModel, SocietyModel, UserModel} from "../../models/user.model";


@Injectable({
  providedIn: 'root'
})
export class AuthService {
  firebaseUser: firebase.User

  constructor(
    private  afAuth: AngularFireAuth,
    private userService: UserService,
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
        this.firebaseUser = res.user;
        res.user.getIdToken().then(token => {
          localStorage.setItem('token', token)
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
          console.log('error custom renew token ' + err);
        })
    })
      .catch(err => {
        console.log('error custom renew token ' + err);
      })
  }

  async register(value) {
    await this.afAuth.createUserWithEmailAndPassword(value.email, value.password)
      .then(res => {
        this.firebaseUser = res.user;
        localStorage.setItem('firebaseUser', JSON.stringify(res.user));
        this.userService.createUser({
          Id: '',
          FirstName: 'Olo',//value.FirstName,
          LastName: 'Chmelo',//value.LastName,
          Email: value.email,
          Uid: this.firebaseUser.uid,
          Avatar: '',
          CreatedAt: new Date()
        }).subscribe(usr => {
          console.log('newUser: ', usr);
          console.log(localStorage.getItem('token'))
          console.log('bol stary a teraz novy')
          this.renewToken();
        });
      })
      .catch(err => {
          console.log(err)
          this.firebaseUser = null;
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
    localStorage.removeItem('user');
    this.firebaseUser = null;
  }

  isLoggedIn(): boolean {
    const user = JSON.parse(localStorage.getItem('user'));
    return user !== null;
  }

  async loginWithGoogle() {
    await this.afAuth.signInWithPopup(new auth.GoogleAuthProvider())
  }

  getToken() {
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
