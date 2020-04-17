import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import {RouterModule} from '@angular/router';

import { AngularFireModule } from '@angular/fire';
import { AngularFirestoreModule } from '@angular/fire/firestore';
import { AngularFireAuthModule } from '@angular/fire/auth';

// import { environment } from '../environments/environment';
import { AgmCoreModule } from '@agm/core';
import { environment } from '../environments/environment';
import {HttpClientModule} from "@angular/common/http";
import { GoogleMapComponent } from './components/google-map/google-map.component';
import {LocationService} from "./services/location/location.service";
import { UserDetailComponent } from './components/user-detail/user-detail.component';
import { LoginComponent } from './components/login/login.component';
import {ReactiveFormsModule} from "@angular/forms";
import { TopBarComponent } from './components/top-bar/top-bar.component';
import { AppRoutingModule } from './app-routing.module';

@NgModule({
  declarations: [
    AppComponent,
    GoogleMapComponent,
    UserDetailComponent,
    LoginComponent,
    TopBarComponent
  ],
  imports: [
    BrowserModule,
    RouterModule,
    AngularFireModule.initializeApp(environment.firebaseConfig),
    AngularFirestoreModule,
    AngularFireAuthModule,
    HttpClientModule,
    AgmCoreModule.forRoot({
      apiKey: environment.googleMapsApiKey,
      libraries: ['places', 'geometry']
    }),
    ReactiveFormsModule,
    AppRoutingModule,
  ],
  providers: [LocationService],
  bootstrap: [AppComponent]
})
export class AppModule { }
