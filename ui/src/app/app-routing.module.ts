import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {UserDetailComponent} from "./components/user-detail/user-detail.component";
import {LoginComponent} from "./components/login/login.component";
import {MyProfileComponent} from "./components/my-profile/my-profile.component";
import {RegisterComponent} from "./components/register/register.component";
import {SocietiesComponent} from "./components/societies/societies/societies.component";

const routes: Routes = [
  { path: '', redirectTo:'me', pathMatch: 'full'},
  { path: 'me', component: MyProfileComponent},
  { path: 'societies', component: SocietiesComponent},
  //{ path: 'profile/tu-daj-user-id', component: UserDetailComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
