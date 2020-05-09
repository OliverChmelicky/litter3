import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppComponent} from './app.component';
import {AgmCoreModule} from '@agm/core';
import {RouterModule} from '@angular/router';

import {AngularFireModule} from '@angular/fire';
import {AngularFirestoreModule} from '@angular/fire/firestore';
import {AngularFireAuthModule} from '@angular/fire/auth';
import {environment} from '../environments/environment';
import {HTTP_INTERCEPTORS, HttpClientModule} from "@angular/common/http";
import {GoogleMapComponent} from './components/google-map/google-map.component';
import {UserDetailComponent} from './components/user-detail/user-detail.component';
import {LoginComponent} from './components/login/login.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { TopBarComponent } from './components/top-bar/top-bar.component';
import { AppRoutingModule } from './app-routing.module';
import { MyProfileComponent } from './components/my-profile/my-profile.component';
import { RegisterComponent } from './components/register/register.component';
import {TokenHeadders} from "./interceptors/token-headders";
import { SocietiesComponent } from './components/societies/societies.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {MatPaginatorModule} from "@angular/material/paginator";
import {MatTableModule} from "@angular/material/table";
import { SocietyDetailsComponent } from './components/society-details/society-details.component';
import { CreateTrashComponent } from './components/create-trash/create-trash.component';
import {MatSliderModule} from "@angular/material/slider";
import {MatSelectModule} from "@angular/material/select";
import {MatCheckboxModule} from "@angular/material/checkbox";
import {MatButtonModule} from "@angular/material/button";
import { NavListComponent } from './components/nav-list/nav-list.component';
import {MatListModule} from "@angular/material/list";
import {MatIconModule} from "@angular/material/icon";
import { TrashDetailsComponent } from './components/trash-details/trash-details.component';
import {MatMenuModule} from "@angular/material/menu";

@NgModule({
  declarations: [
    AppComponent,
    GoogleMapComponent,
    UserDetailComponent,
    LoginComponent,
    TopBarComponent,
    MyProfileComponent,
    RegisterComponent,
    SocietiesComponent,
    SocietyDetailsComponent,
    CreateTrashComponent,
    NavListComponent,
    TrashDetailsComponent,
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
    FormsModule,
    BrowserAnimationsModule,
    MatPaginatorModule,
    MatTableModule,
    MatSliderModule,
    MatSelectModule,
    MatCheckboxModule,
    MatButtonModule,
    MatListModule,
    MatIconModule,
    MatMenuModule,
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: TokenHeadders, multi: true },
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
