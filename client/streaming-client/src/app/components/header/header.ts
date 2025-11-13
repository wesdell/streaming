import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';

import { AuthService } from '../../services/auth';

@Component({
  selector: 'header',
  standalone: true,
  templateUrl: './header.html'
})
export class Header {
  private router = inject(Router);
  authService = inject(AuthService);

  get auth() {
    return this.authService.auth;
  }

  goToLogin() {
    this.router.navigate(['/login']);
  }

  goToRegister() {
    this.router.navigate(['/register']);
  }

  handleLogout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
