// Source https://www.youtube.com/watch?v=qP5zw7fjQgo

import {Injectable} from '@angular/core';

import {auth} from 'firebase/app';
import {AngularFireAuth} from '@angular/fire/auth';
import {
  AngularFirestore,
  AngularFirestoreDocument
} from '@angular/fire/firestore';

import {Observable, of} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {UserModel} from '../../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  user$: Observable<UserModel>;

  constructor(
    private afAuth: AngularFireAuth,
    private afs: AngularFirestore,
  ) {

  }

  async googleSignin() {
    const provider = new auth.GoogleAuthProvider();
    const credential = await this.afAuth.auth.signInWithPopup(provider);
    return this.updateUserData(credential.user);
  }

  async signOut() {
    await this.afAuth.auth.signOut();
  }

  private updateUserData({uid, email, displayName, photoURL}: UserModel) {
    // Sets user data to firestore on login
    const userRef: AngularFirestoreDocument<UserModel> = this.afs.doc(`users/${uid}`);

    const data = {
      uid,
      email,
      displayName,
      photoURL
    };

    return userRef.set(data, {merge: true});
  }

  RegisterPasswdAuth() {
    this.afAuth.auth.createUserWithEmailAndPassword("idem@plavat.tu", "test11")
      .catch(function (error) {
        console.log(error.message)
      });
  }

  LoginPasswdAuth() {
    this.afAuth.auth.signInWithEmailAndPassword("idem@plavat.tu", "test11").then(data => console.log(data.user.getIdToken()))
      .catch(function (error) {
        console.log(error.message)
      });
  }

}
