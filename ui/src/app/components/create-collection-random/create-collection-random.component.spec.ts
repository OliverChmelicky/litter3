import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateCollectionRandomComponent } from './create-collection-random.component';

describe('CreateCollectionRandomComponent', () => {
  let component: CreateCollectionRandomComponent;
  let fixture: ComponentFixture<CreateCollectionRandomComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateCollectionRandomComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateCollectionRandomComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
