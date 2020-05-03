import { TestBed } from '@angular/core/testing';

import { SocietyService } from './society.service';

describe('SocietyService', () => {
  let service: SocietyService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SocietyService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
