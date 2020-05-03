import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {UserDetailComponent} from "./components/user-detail/user-detail.component";
import {LoginComponent} from "./components/login/login.component";
import {MyProfileComponent} from "./components/my-profile/my-profile.component";
import {RegisterComponent} from "./components/register/register.component";
import {SocietiesComponent} from "./components/societies/societies.component";
import {SocietyDetailsComponent} from "./components/society-details/society-details.component";
import {GoogleMapComponent} from "./components/google-map/google-map.component";

const routes: Routes = [
  { path: '', redirectTo:'map', pathMatch: 'full'},
  { path: 'map', component: GoogleMapComponent},
  { path: 'me', component: MyProfileComponent},
  { path: 'societies', component: SocietiesComponent},
  //{ path: 'profile/tu-daj-user-id', component: UserDetailComponent },
  {path: 'societies/:societyId', component: SocietyDetailsComponent},
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
