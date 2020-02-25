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
import { GoogleMapComponent } from './google-map/google-map.component';
import {LocationService} from "./services/location/location.service";

@NgModule({
  declarations: [
    AppComponent,
    GoogleMapComponent
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
  ],
  providers: [LocationService],
  bootstrap: [AppComponent]
})
export class AppModule { }
